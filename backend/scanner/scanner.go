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

func (s *Scanner) ScanDirectories() error {
	s.setRunning()
	defer s.setFinished()
	directoryItems, err := s.collectDirectoryItems()
	if err != nil {
		return fmt.Errorf("collecting directory items: %w", err)
	}

	s.Status = Status{}
	start := time.Now()

	log.Printf("starting import at %v", start)
	s.Status.CountTotal = len(directoryItems)
	s.Updates <- struct{}{}

	defer func() {
		log.Printf("finished import in %v", time.Since(start))
		s.Status.CountProcessed = len(directoryItems)
		s.Updates <- struct{}{}
	}()

	for idx, item := range directoryItems {
		hash, err := s.scanDirectoryItem(item)
		if err != nil {
			s.Status.AddError(err)
			s.Updates <- struct{}{}
			continue
		}
		if hash == "" {
			continue
		}
		s.Status.LastHash = hash
		s.Status.CountProcessed = idx + 1
		s.Updates <- struct{}{}
	}
	return nil
}

type dirItem struct {
	dirAlias string
	dir      string
	fileName string
	modTime  time.Time
}

func (s *Scanner) collectDirectoryItems() ([]*dirItem, error) {
	var items []*dirItem
	for alias, dir := range s.Directories {
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
			modTime := info.ModTime()
			items = append(items, &dirItem{
				alias, dir, fileName, modTime,
			})
		}
	}
	return items, nil
}

func (s *Scanner) scanDirectoryItem(item *dirItem) (string, error) {
	_, err := s.DB.GetDirInfo(context.Background(), item.dirAlias, item.fileName)
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
	if err := s.Importer.ImportScreenshot(decoded, timestamp, item.dirAlias, item.fileName); err != nil {
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
