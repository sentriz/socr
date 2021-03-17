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
	"github.com/fsnotify/fsnotify"
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
		hash, err := s.scanDirectoryItem(item.dirAlias, item.dir, item.fileName, item.modTime)
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

func (s *Scanner) WatchUpdates() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	for alias, dir := range s.Directories {
		if err = watcher.Add(dir); err != nil {
			return fmt.Errorf("add watcher for %q: %w", alias, err)
		}
		log.Printf("starting watcher for %q", dir)
	}
	for event := range watcher.Events {
		if event.Op&fsnotify.Create != fsnotify.Create {
			continue
		}
		if strings.HasSuffix(event.Name, ".tmp") {
			continue
		}
		dir := filepath.Dir(event.Name)
		dirAlias, ok := dirAliasFromDir(s.Directories, dir)
		if !ok {
			continue
		}
		fileName := filepath.Base(event.Name)
		modTime := time.Now()
		if _, err := s.scanDirectoryItem(dirAlias, dir, fileName, modTime); err != nil {
			log.Printf("error scanning directory item with event %v: %v", event, err)
		}
	}
	return nil
}

type collectedDirectoryItem struct {
	dirAlias string
	dir      string
	fileName string
	modTime  time.Time
}

func (s *Scanner) collectDirectoryItems() ([]*collectedDirectoryItem, error) {
	var items []*collectedDirectoryItem
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
			items = append(items, &collectedDirectoryItem{
				alias, dir, fileName, modTime,
			})
		}
	}
	return items, nil
}

func (s *Scanner) scanDirectoryItem(dirAlias, dir, fileName string, modTime time.Time) (string, error) {
	_, err := s.DB.GetDirInfo(context.Background(), dirAlias, fileName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("getting dir info: %w", err)
	}
	if err == nil {
		return "", nil
	}

	log.Printf("importing new screenshot. alias %q, filename %q", dirAlias, fileName)

	filePath := filepath.Join(dir, fileName)
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %v", err)
	}

	decoded, err := importer.DecodeImage(raw)
	if err != nil {
		return "", fmt.Errorf("decode screenshot: %v", err)
	}

	timestamp := guessFileCreated(fileName, modTime)
	if err := s.Importer.ImportScreenshot(decoded, timestamp, dirAlias, fileName); err != nil {
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

func dirAliasFromDir(directories map[string]string, dir string) (string, bool) {
	for k, v := range directories {
		if v == dir {
			return k, true
		}
	}
	return "", false
}
