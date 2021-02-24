package importer

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
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
	isRunning   *int32
	DB          *db.Conn
	Hasher      hasher.Hasher
	Directories map[string]string
}

type StatusError struct {
	Time  time.Time `json:"time"`
	Error string    `json:"error"`
}

type Status struct {
	Errors         []StatusError `json:"errors,omitempty"`
	LastID         string        `json:"last_id,omitempty"`
	CountProcessed int           `json:"count_processed"`
	CountTotal     int           `json:"count_total"`
}

func (s *Status) AddError(err error) {
	s.Errors = append(s.Errors, StatusError{
		Time:  time.Now(),
		Error: err.Error(),
	})
}

func (i *Importer) IsRunning() bool { return atomic.LoadInt32(i.isRunning) == 1 }
func (i *Importer) SetRunning()     { atomic.StoreInt32(i.isRunning, 1) }
func (i *Importer) SetFinished()    { atomic.StoreInt32(i.isRunning, 0) }

func (i *Importer) ScanDirectories() error {
	start := time.Now()

	for alias, dir := range i.Directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("listing dir: %w", err)
		}
		for _, file := range files {
			filename := file.Name()
			screenshot, err := i.scanDirectoryItem(alias, dir, file)
			if err != nil {
				return fmt.Errorf("import or skip %q: %w", filename, err)
			}
			log.Printf("imported screenshot. id %d filename %q", screenshot.ID, filename)
		}
	}

	log.Printf("finished import in %v", time.Since(start))
	return nil
}

func (i *Importer) scanDirectoryItem(alias, dir string, file fs.FileInfo) (*db.Screenshot, error) {
	screenshot, err := i.DB.GetScreenshotByPath(context.Background(), db.GetScreenshotByPathParams{
		DirectoryAlias: alias,
		Filename:       file.Name(),
	})
	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("getting screenshot by path: %v", err)
	case err == nil:
		return &screenshot, nil
	}

	fileName := file.Name()
	filePath := filepath.Join(dir, fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading from disk: %v", err)
	}

	hash, err := i.Hasher.Hash(bytes)
	if err != nil {
		return nil, fmt.Errorf("hashing screenshot: %v", err)
	}

	timestamp := guessFileCreated(file)
	imported, err := i.ImportScreenshot(hash, timestamp, alias, fileName, bytes)
	if err != nil {
		return nil, fmt.Errorf("importing screenshot: %v", err)
	}

	return imported, nil
}

func (i *Importer) ImportScreenshot(id uint64, timestamp time.Time, dirAlias string, filename string, raw []byte) (*db.Screenshot, error) {
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

	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := imagery.FormatPNG.Encode(imageEncoded, imageBig); err != nil {
		return nil, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	blocks, err := imagery.ExtractText(imageEncoded.Bytes())
	if err != nil {
		return nil, fmt.Errorf("extract image text: %w", err)
	}

	size := image.Bounds().Size()
	screenshotArgs := db.CreateScreenshotParams{
		ID:             int64(id),
		Timestamp:      timestamp,
		DirectoryAlias: dirAlias,
		Filename:       filename,
		DimWidth:       int32(size.X),
		DimHeight:      int32(size.Y),
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	}

	screenshost, err := i.DB.CreateScreenshot(context.Background(), screenshotArgs)
	if err != nil {
		return nil, fmt.Errorf("inserting screenshot: %w", err)
	}

	tx, err := i.DB.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	q := i.DB.WithTx(tx)
	for idx, block := range blocks {
		rect := imagery.ScaleDownRect(block.Box)
		err := q.CreateBlock(context.Background(), db.CreateBlockParams{
			ScreenshotID: screenshost.ID,
			Index:        int16(idx),
			MinX:         int16(rect[0]),
			MinY:         int16(rect[1]),
			MaxX:         int16(rect[2]),
			MaxY:         int16(rect[3]),
			Body:         block.Word,
		})
		if err != nil {
			return nil, fmt.Errorf("inserting block: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("end transaction: %w", err)
	}

	return &screenshost, nil
}

func guessFileCreated(file os.FileInfo) time.Time {
	filename := file.Name()
	filename = strings.ToLower(filename)
	filename = strings.TrimPrefix(filename, "img_")
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = strings.ReplaceAll(filename, "_", "")

	guessed, err := dateparse.ParseLocal(filename)
	if err != nil {
		log.Printf("couldn't guess timestamp of %q, using mod time", file.Name())
		return file.ModTime()
	}

	return guessed
}
