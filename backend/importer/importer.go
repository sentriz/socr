package importer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/araddon/dateparse"
	"github.com/blevesearch/bleve"

	"go.senan.xyz/socr/db"
	"go.senan.xyz/socr/hasher"
	"go.senan.xyz/socr/screenshot"
)

type Importer struct {
	isRunning   *int32
	DB          db.DB
	Hasher      hasher.Hasher
	Index       bleve.Index
	Directories screenshot.Directories
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

func (i *Importer) Scan() error {
	for _, dir := range i.Directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("listing dir: %w", err)
		}

		for _, file := range files {
			timeLast, err := i.DB.GetModTime(dir, file.Name())
			if err != nil {
				return fmt.Errorf("get last mod time: %v", err)
			}
			if timeLast != nil {
				continue
			}

			bytes, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return fmt.Errorf("reading from disk: %v", err)
			}

			hash, err := i.Hasher.Hash(bytes)
			if err != nil {
				return fmt.Errorf("hashing screenshot: %v", err)
			}

			screenshotID := i.Hasher.Format(hash)
			timestamp := guessFileCreated(file)
			screenshot, err := screenshot.FromBytesWithHash(bytes, screenshotID, timestamp)
			if err != nil {
				return fmt.Errorf("processing screenshot: %v", err)
			}

			if err := i.Index.Index(screenshot.ID, screenshot); err != nil {
				return fmt.Errorf("indexing screenshot: %w", err)
			}
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
