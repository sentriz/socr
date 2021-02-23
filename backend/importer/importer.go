package importer

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (i *Importer) Import() error {
	for alias, dir := range i.Directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("listing dir: %w", err)
		}

		for _, file := range files {
			_, err := i.DB.GetScreenshotByPath(context.Background(), db.GetScreenshotByPathParams{
				DirectoryAlias: alias,
				Filename:       file.Name(),
			})
			switch {
			case err != nil && !errors.Is(err, sql.ErrNoRows):
				return fmt.Errorf("getting screenshot by path: %v", err)
			case err == nil:
				continue
			}

			filename := file.Name()
			bytes, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("reading from disk: %v", err)
			}

			hash, err := i.Hasher.Hash(bytes)
			if err != nil {
				return fmt.Errorf("hashing screenshot: %v", err)
			}

			timestamp := guessFileCreated(file)
			screenshot, err := i.importScreenshot(hash, timestamp, alias, filename, bytes)
			if err != nil {
				return fmt.Errorf("importing screenshot: %v", err)
			}

			log.Printf("imported screenshot. id %q filename %q", screenshot.ID, screenshot.Filename)
		}
	}

	return nil
}

func guessFileCreated(file os.FileInfo) time.Time {
	filename := file.Name()
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = strings.ReplaceAll(filename, "_", "")

	guessed, err := dateparse.ParseLocal(filename)
	if err != nil {
		log.Printf("couldn't guess timestamp of %q, using mod time", file.Name())
		return file.ModTime()
	}

	return guessed
}

func (i *Importer) importScreenshot(id uint64, timestamp time.Time, dirAlias string, filename string, raw []byte) (*db.Screenshot, error) {
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

	screenshotArgs := db.CreateScreenshotParams{
		ID:             int64(id),
		Timestamp:      timestamp,
		DirectoryAlias: dirAlias,
		Filename:       filename,
		DimWidth:       0,
		DimHeight:      0,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	}

	screenshost, err := i.DB.CreateScreenshot(context.Background(), screenshotArgs)
	if err != nil {
		return nil, fmt.Errorf("inserting screenshot: %w", err)
	}

	for _, block := range blocks {
		rect := imagery.ScaleDownRect(block.Box)
		err := i.DB.CreateBlock(context.Background(), db.CreateBlockParams{
			MinX: int16(rect[0]), MinY: int16(rect[1]),
			MaxX: int16(rect[2]), MaxY: int16(rect[3]),
			Body: block.Word,
		})
		if err != nil {
			return nil, fmt.Errorf("inserting block: %w", err)
		}
	}

	return &screenshost, nil
}
