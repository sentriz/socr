package importer

import (
	"bytes"
	"context"
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
	UpdatesScreenshot     chan string
}

type StatusError struct {
	Time  time.Time `json:"time"`
	Error string    `json:"error"`
}

type Status struct {
	Errors         []StatusError `json:"errors"`
	LastHash       string        `json:"last_hash"`
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
		hash, err := i.scanDirectoryItem(item)
		if err != nil {
			i.Status.AddError(err)
			i.UpdatesScan <- struct{}{}
			continue
		}

		i.Status.LastHash = hash
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

func (i *Importer) scanDirectoryItem(item *collected) (string, error) {
	row, err := i.DB.GetScreenshotByPath(context.Background(), item.dirAlias, item.fileName)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return "", fmt.Errorf("getting screenshot by path: %v", err)
	case err == nil:
		return row.Hash, nil
	}

	log.Printf("importing new screenshot. alias %q, filename %q", item.dirAlias, item.fileName)

	filePath := filepath.Join(item.dir, item.fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading from disk: %v", err)
	}

	hash := hasher.Hash(bytes)
	timestamp := guessFileCreated(item.fileName, item.modTime)
	if err := i.ImportScreenshot(hash, timestamp, item.dirAlias, item.fileName, bytes); err != nil {
		return "", fmt.Errorf("importing screenshot: %v", err)
	}

	return hash, nil
}

func (i *Importer) ImportScreenshot(hash string, timestamp time.Time, dirAlias, fileName string, raw []byte) error {
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
	id, err := i.importScreenshotProperties(hash, image, timestamp, dirAlias, fileName)
	if err != nil {
		return fmt.Errorf("import screenshot props: %w", err)
	}

	// insert to blocks
	if err := i.importScreenshotBlocks(id, image); err != nil {
		return fmt.Errorf("import screenshot props: %w", err)
	}

	i.UpdatesScreenshot <- hash

	return nil
}

func (i *Importer) importScreenshotProperties(hash string, image image.Image, timestamp time.Time, dirAlias, filename string) (int, error) {
	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, fmt.Errorf("calculate blurhash: %w", err)
	}

	size := image.Bounds().Size()
	screenshotArgs := db.CreateScreenshotParams{
		Hash:           hash,
		Timestamp:      timestamp,
		DirectoryAlias: dirAlias,
		Filename:       filename,
		DimWidth:       size.X,
		DimHeight:      size.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	}

	row, err := i.DB.CreateScreenshot(context.Background(), screenshotArgs)
	if err != nil {
		return 0, fmt.Errorf("inserting screenshot: %w", err)
	}

	return row.ID, nil
}

func (i *Importer) importScreenshotBlocks(id int, image image.Image) error {
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
			ScreenshotID: id,
			Index:        idx,
			MinX:         rect[0],
			MinY:         rect[1],
			MaxX:         rect[2],
			MaxY:         rect[3],
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
