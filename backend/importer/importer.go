package importer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	Running           *int32
	DB                *db.DB
	Directories       map[string]string
	Status            Status
	UpdatesScan       chan struct{}
	UpdatesScreenshot chan string
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

	directoryItems, err := i.collectDirectoryItems()
	if err != nil {
		return fmt.Errorf("collecting directory items: %w", err)
	}

	i.Status = Status{}
	i.Status.CountTotal = len(directoryItems)
	i.UpdatesScan <- struct{}{}

	start := time.Now()
	log.Printf("starting import at %v", start)

	for idx, item := range directoryItems {
		i.Status.CountProcessed = idx + 1
		hash, err := i.scanDirectoryItem(item)
		if err != nil {
			i.Status.AddError(err)
			i.UpdatesScan <- struct{}{}
			continue
		}
		if hash == "" {
			continue
		}
		i.Status.LastHash = hash
		i.UpdatesScan <- struct{}{}
	}

	i.UpdatesScan <- struct{}{}
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
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("listing dir %q: %w", dir, err)
		}
		for _, file := range files {
			fileName := file.Name()
			info, err := file.Info()
			if err != nil {
				return nil, fmt.Errorf("get file info %q: %w", fileName, err)
			}
			items = append(items, &collected{
				dirAlias: alias,
				dir:      dir,
				fileName: fileName,
				modTime:  info.ModTime(),
			})
		}
	}
	return items, nil
}

func (i *Importer) scanDirectoryItem(item *collected) (string, error) {
	_, err := i.DB.GetDirInfo(context.Background(), item.dirAlias, item.fileName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("getting dir info: %w", err)
	}
	if err == nil {
		return "", nil
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
	// insert to screenshots and blocks
	id, err := i.importScreenshot(hash, timestamp, raw)
	if err != nil {
		return fmt.Errorf("props and blocks: %w", err)
	}

	// insert to dir info
	if err := i.importScreenshotDirInfo(id, dirAlias, fileName); err != nil {
		return fmt.Errorf("dir info: %w", err)
	}

	i.UpdatesScreenshot <- hash

	return nil
}

func (i *Importer) importScreenshot(hash string, timestamp time.Time, raw []byte) (int, error) {
	row, err := i.DB.GetScreenshotByHash(context.Background(), hash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("getting screenshot by hash: %w", err)
	}
	if row.ID != 0 {
		// we already have this screenshot
		return row.ID, nil
	}

	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return 0, fmt.Errorf("unrecognised format: %s", mime)
	}
	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return 0, fmt.Errorf("decoding: %s", mime)
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, fmt.Errorf("calculate blurhash: %w", err)
	}

	propSize := image.Bounds().Size()
	screenshotRow, err := i.DB.CreateScreenshot(context.Background(), db.CreateScreenshotParams{
		Hash:           hash,
		Timestamp:      timestamp,
		DimWidth:       propSize.X,
		DimHeight:      propSize.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if err != nil {
		return 0, fmt.Errorf("inserting screenshot: %w", err)
	}

	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := imagery.FormatPNG.Encode(imageEncoded, imageBig); err != nil {
		return 0, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	blocks, err := imagery.ExtractText(imageEncoded.Bytes())
	if err != nil {
		return 0, fmt.Errorf("extract image text: %w", err)
	}

	batch := &pgx.Batch{}
	for idx, block := range blocks {
		rect := imagery.ScaleDownRect(block.Box)
		i.DB.CreateBlockBatch(batch, db.CreateBlockParams{
			ScreenshotID: screenshotRow.ID,
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
		return 0, fmt.Errorf("end transaction: %w", err)
	}

	return screenshotRow.ID, nil
}

func (i *Importer) importScreenshotDirInfo(id int, dirAlias string, fileName string) error {
	_, err := i.DB.CreateDirInfo(context.Background(), db.CreateDirInfoParams{
		ScreenshotID:   id,
		Filename:       fileName,
		DirectoryAlias: dirAlias,
	})
	if err != nil {
		return fmt.Errorf("insert info dir infos: %w", err)
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
