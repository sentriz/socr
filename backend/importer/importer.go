package importer

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
	"go.senan.xyz/socr/backend/imagery"
)

type Importer struct {
	Running               *int32
	DB                    *db.DB
	Directories           map[string]string
	DirectoriesUploadsKey string
	Status                Status
	UpdatesScan           chan struct{}
	UpdatesScreenshot     chan int64
}

type StatusError struct {
	Time  time.Time
	Error string
}

type Status struct {
	Errors         []StatusError
	LastID         int64
	CountProcessed int
	CountTotal     int
}

func (s *Status) AddError(err error) {
	s.Errors = append(s.Errors, StatusError{
		Time:  time.Now(),
		Error: err.Error(),
	})
}

func (i *Importer) IsRunning() bool { return atomic.LoadInt32(i.Running) == 1 }
func (i *Importer) setRunning()     { atomic.StoreInt32(i.Running, 1) }
func (i *Importer) setFinished()    { atomic.StoreInt32(i.Running, 0) }

func (i *Importer) ScanDirectories() error {
	i.setRunning()
	defer i.setFinished()

	i.Status = Status{}
	i.UpdatesScan <- struct{}{}

	directoryItems, err := i.collectDirectoryItems()
	if err != nil {
		return fmt.Errorf("collecting directory items: %w", err)
	}

	start := time.Now()
	log.Printf("starting import at %v", start)

	for _, item := range directoryItems {
		id, err := i.scanDirectoryItem(item)
		if err != nil {
			i.Status.AddError(err)
			i.UpdatesScan <- struct{}{}
			continue
		}

		i.Status.LastID = id
		i.Status.CountProcessed++
		i.Status.CountTotal = len(directoryItems)
		i.UpdatesScan <- struct{}{}
	}

	log.Printf("finished import in %v", time.Since(start))
	return nil
}

type collected struct {
	dirAlias string
	dir      string
	fileName string
	modTime  time.Time
}

func (i *Importer) collectDirectoryItems() ([]*collected, error) {
	var items []*collected
	for alias, dir := range i.Directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("listing dir: %w", err)
		}
		for _, file := range files {
			items = append(items, &collected{
				dirAlias: alias,
				dir:      dir,
				fileName: file.Name(),
			})
		}
	}
	return items, nil
}

func (i *Importer) scanDirectoryItem(item *collected) (int64, error) {
	row, err := i.DB.GetScreenshotByPath(context.Background(), item.dirAlias, item.fileName)
	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return 0, fmt.Errorf("getting screenshot by path: %v", err)
	case err == nil:
		return row.ID.Int, nil
	}

	log.Printf("importing new screenshot. alias %q, filename %q", item.dirAlias, item.fileName)

	filePath := filepath.Join(item.dir, item.fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("reading from disk: %v", err)
	}

	id, err := hasher.Hash(bytes)
	if err != nil {
		return 0, fmt.Errorf("hashing screenshot: %v", err)
	}

	timestamp := guessFileCreated(item.fileName, item.modTime)
	if err := i.ImportScreenshot(id, timestamp, item.dirAlias, item.fileName, bytes); err != nil {
		return 0, fmt.Errorf("importing screenshot: %v", err)
	}

	return id, nil
}

func (i *Importer) ImportScreenshot(id int64, timestamp time.Time, dirAlias, fileName string, raw []byte) error {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return fmt.Errorf("unrecognised format: %s", mime)
	}

	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return fmt.Errorf("decoding: %s", mime)
	}

	// insert to screenshots
	if err := i.importScreenshotProperties(id, image, timestamp, dirAlias, fileName); err != nil {
		return fmt.Errorf("import screenshot props: %w", err)
	}

	// insert to blocks
	if err := i.importScreenshotBlocks(id, image); err != nil {
		return fmt.Errorf("import screenshot props: %w", err)
	}

	i.UpdatesScreenshot <- id

	return nil
}

func (i *Importer) importScreenshotProperties(id int64, image image.Image, timestamp time.Time, dirAlias, filename string) error {
	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return fmt.Errorf("calculate blurhash: %w", err)
	}

	size := image.Bounds().Size()
	screenshotArgs := db.CreateScreenshotParams{
		ID:             int(id),
		Timestamp:      pgtype.Timestamp{Time: timestamp},
		DirectoryAlias: dirAlias,
		Filename:       filename,
		DimWidth:       int32(size.X),
		DimHeight:      int32(size.Y),
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	}

	if _, err := i.DB.CreateScreenshot(context.Background(), screenshotArgs); err != nil {
		return fmt.Errorf("inserting screenshot: %w", err)
	}

	return nil
}

func (i *Importer) importScreenshotBlocks(screenshotID int64, image image.Image) error {
	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := imagery.FormatPNG.Encode(imageEncoded, imageBig); err != nil {
		return fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	blocks, err := imagery.ExtractText(imageEncoded.Bytes())
	if err != nil {
		return fmt.Errorf("extract image text: %w", err)
	}

	batch := &pgx.Batch{}
	for idx, block := range blocks {
		rect := imagery.ScaleDownRect(block.Box)
		i.DB.CreateBlockBatch(batch, db.CreateBlockParams{
			ScreenshotID: int(screenshotID),
			Index:        int16(idx),
			MinX:         int16(rect[0]),
			MinY:         int16(rect[1]),
			MaxX:         int16(rect[2]),
			MaxY:         int16(rect[3]),
			Body:         block.Word,
		})
	}

	results := i.DB.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return fmt.Errorf("end transaction: %w", err)
	}

	return nil
}

func guessFileCreated(fileName string, modTime time.Time) time.Time {
	fileName = strings.ToLower(fileName)
	fileName = strings.TrimPrefix(fileName, "img_")
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileName = strings.ReplaceAll(fileName, "_", "")

	guessed, err := dateparse.ParseLocal(fileName)
	if err != nil {
		log.Printf("couldn't guess timestamp of %q, using mod time", fileName)
		return modTime
	}

	return guessed
}
