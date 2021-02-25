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

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
	"go.senan.xyz/socr/backend/imagery"
)

type Importer struct {
	Running               *int32
	DB                    *db.Conn
	Directories           map[string]string
	DirectoriesUploadsKey string
	Status                Status
	UpdatesScan           chan struct{}
	UpdatesScreenshot     chan *db.Screenshot
}

type StatusError struct {
	Time  time.Time `json:"time"`
	Error string    `json:"error"`
}

type Status struct {
	Errors         []StatusError `json:"errors,omitempty"`
	LastID         hasher.ID     `json:"last_id,omitempty"`
	CountProcessed int           `json:"count_processed"`
	CountTotal     int           `json:"count_total"`
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
		screenshot, err := i.scanDirectoryItem(item)
		if screenshot == nil {
			continue
		}
		if err != nil {
			i.Status.AddError(err)
			i.UpdatesScan <- struct{}{}
			continue
		}

		i.Status.LastID = screenshot.ID
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

func (i *Importer) scanDirectoryItem(item *collected) (*db.Screenshot, error) {
	screenshot, err := i.DB.GetScreenshotByPath(context.Background(), db.GetScreenshotByPathParams{
		DirectoryAlias: item.dirAlias,
		Filename:       item.fileName,
	})
	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("getting screenshot by path: %v", err)
	case err == nil:
		return &screenshot, nil
	}

	log.Printf("importing new screenshot. alias %q, filename %q", item.dirAlias, item.fileName)

	filePath := filepath.Join(item.dir, item.fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading from disk: %v", err)
	}

	hash, err := hasher.Hash(bytes)
	if err != nil {
		return nil, fmt.Errorf("hashing screenshot: %v", err)
	}

	timestamp := guessFileCreated(item.fileName, item.modTime)
	imported, err := i.ImportScreenshot(hash, timestamp, item.dirAlias, item.fileName, bytes)
	if err != nil {
		return nil, fmt.Errorf("importing screenshot: %v", err)
	}

	return imported, nil
}

func (i *Importer) ImportScreenshot(id hasher.ID, timestamp time.Time, dirAlias, fileName string, raw []byte) (*db.Screenshot, error) {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unrecognised format: %s", mime)
	}

	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return nil, fmt.Errorf("decoding: %s", mime)
	}

	// insert to screenshots
	screenshot, err := i.importScreenshotProperties(id, image, timestamp, dirAlias, fileName)
	if err != nil {
		return nil, fmt.Errorf("import screenshot props: %w", err)
	}

	// insert to blocks
	err = i.importScreenshotBlocks(id, image)
	if err != nil {
		return nil, fmt.Errorf("import screenshot props: %w", err)
	}

	i.UpdatesScreenshot <- screenshot

	return screenshot, nil
}

func (i *Importer) importScreenshotProperties(id hasher.ID, image image.Image, timestamp time.Time, dirAlias, filename string) (*db.Screenshot, error) {
	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	size := image.Bounds().Size()
	screenshotArgs := db.CreateScreenshotParams{
		ID:             id,
		Timestamp:      timestamp,
		DirectoryAlias: dirAlias,
		Filename:       filename,
		DimWidth:       int32(size.X),
		DimHeight:      int32(size.Y),
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	}

	screenshot, err := i.DB.CreateScreenshot(context.Background(), screenshotArgs)
	if err != nil {
		return nil, fmt.Errorf("inserting screenshot: %w", err)
	}

	return &screenshot, nil
}

func (i *Importer) importScreenshotBlocks(screenshotID hasher.ID, image image.Image) error {
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

	tx, err := i.DB.Conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	q := i.DB.WithTx(tx)
	for idx, block := range blocks {
		rect := imagery.ScaleDownRect(block.Box)
		err := q.CreateBlock(context.Background(), db.CreateBlockParams{
			ScreenshotID: screenshotID,
			Index:        int16(idx),
			MinX:         int16(rect[0]),
			MinY:         int16(rect[1]),
			MaxX:         int16(rect[2]),
			MaxY:         int16(rect[3]),
			Body:         block.Word,
		})
		if err != nil {
			return fmt.Errorf("inserting block: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
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
