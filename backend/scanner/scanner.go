package scanner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jackc/pgx/v4"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/importer"
)

type Scanner struct {
	Running     *int32
	DB          *db.DB
	Directories map[string]string
	Status      Status
	Importer    *importer.Importer
	Updates     chan struct{}
}

func (i *Scanner) ScanDirectories() error {
	i.setRunning()
	defer i.setFinished()

	directoryItems, err := i.collectDirectoryItems()
	if err != nil {
		return fmt.Errorf("collecting directory items: %w", err)
	}

	i.Status = Status{}
	i.Status.CountTotal = len(directoryItems)
	i.Updates <- struct{}{}

	start := time.Now()
	log.Printf("starting import at %v", start)

	for idx, item := range directoryItems {
		i.Status.CountProcessed = idx + 1

		hash, err := i.scanDirectoryItem(item)
		if err != nil {
			i.Status.AddError(err)
			i.Updates <- struct{}{}
			continue
		}
		if hash == "" {
			continue
		}
		i.Status.LastHash = hash
		i.Updates <- struct{}{}
	}

	i.Updates <- struct{}{}
	log.Printf("finished import in %v", time.Since(start))
	return nil
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
	if len(s.Errors) > 20 {
		s.Errors = s.Errors[1:]
	}
}

func (s *Scanner) IsRunning() bool { return atomic.LoadInt32(s.Running) == 1 }
func (s *Scanner) setRunning()     { atomic.StoreInt32(s.Running, 1) }
func (s *Scanner) setFinished()    { atomic.StoreInt32(s.Running, 0) }

type dirItem struct {
	dirAlias string
	dir      string
	fileName string
	modTime  time.Time
}

func (i *Scanner) collectDirectoryItems() ([]*dirItem, error) {
	var items []*dirItem
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
			items = append(items, &dirItem{
				dirAlias: alias,
				dir:      dir,
				fileName: fileName,
				modTime:  info.ModTime(),
			})
		}
	}
	return items, nil
}

func (i *Scanner) scanDirectoryItem(item *dirItem) (string, error) {
	_, err := i.DB.GetDirInfo(context.Background(), item.dirAlias, item.fileName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("getting dir info: %w", err)
	}
	if err == nil {
		return "", nil
	}

	log.Printf("importing new screenshot. alias %q, filename %q", item.dirAlias, item.fileName)

	filePath := filepath.Join(item.dir, item.fileName)
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %v", err)
	}

	decoded, err := importer.DecodeImage(raw)
	if err != nil {
		return "", fmt.Errorf("decode screenshot: %v", err)
	}

	timestamp := guessFileCreated(item.fileName, item.modTime)
	if err := i.Importer.ImportScreenshot(decoded, timestamp, item.dirAlias, item.fileName); err != nil {
		return "", fmt.Errorf("importing screenshot: %v", err)
	}

	return decoded.Hash, nil
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
