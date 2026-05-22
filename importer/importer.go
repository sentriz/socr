package importer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fsnotify/fsnotify"
	"github.com/jackc/pgx/v4"

	"go.senan.xyz/socr/db"
	"go.senan.xyz/socr/directories"
	"go.senan.xyz/socr/imagery"
)

type EncodeFunc func(io.Writer, image.Image) error

type NotifyMediaFunc func(ctx context.Context, hash string)
type NotifyProgressFunc func(ctx context.Context)

type Importer struct {
	db                      *db.DB
	defaultEncoder          EncodeFunc
	defaultMIME             string
	directories             directories.Directories
	directoriesUploadsAlias string
	thumbnailWidth          uint

	statusMu            sync.RWMutex
	status              Status
	jobs                chan *mediaFile
	scanTriggers        chan struct{}
	notifyMediaFuncs    []NotifyMediaFunc
	notifyProgressFuncs []NotifyProgressFunc
}

func New(
	db *db.DB, defaultEncoder EncodeFunc, defaultMIME string,
	directories directories.Directories, directoriesUploadsAlias string, thumbnailWidth uint,
) *Importer {
	return &Importer{
		db:                      db,
		defaultEncoder:          defaultEncoder,
		defaultMIME:             defaultMIME,
		directories:             directories,
		directoriesUploadsAlias: directoriesUploadsAlias,
		thumbnailWidth:          thumbnailWidth,

		jobs:         make(chan *mediaFile),
		scanTriggers: make(chan struct{}, 1),
	}
}

func (i *Importer) TriggerScan() {
	select {
	case i.scanTriggers <- struct{}{}:
	default:
	}
}

func (i *Importer) EnqueueFile(ctx context.Context, alias, dir, fileName string, modTime time.Time) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case i.jobs <- &mediaFile{alias, dir, fileName, modTime}:
		return nil
	}
}

func (i *Importer) RunScanLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-i.scanTriggers:
			if err := i.scanDirectories(ctx); err != nil {
				log.Printf("error scanning directories: %v", err)
			}
		}
	}
}

func (i *Importer) RunPeriodicScan(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			i.TriggerScan()
		}
	}
}

func (i *Importer) AddNotifyMediaFunc(f NotifyMediaFunc) {
	i.notifyMediaFuncs = append(i.notifyMediaFuncs, f)
}

func (i *Importer) AddNotifyProgressFunc(f NotifyProgressFunc) {
	i.notifyProgressFuncs = append(i.notifyProgressFuncs, f)
}

func (i *Importer) Status() Status {
	i.statusMu.RLock()
	defer i.statusMu.RUnlock()
	return i.status
}

func (i *Importer) StartWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case j, ok := <-i.jobs:
			if !ok {
				return nil
			}
			hash, err := i.importMediaFromFile(ctx, j.dirAlias, j.dir, j.fileName, j.modTime)
			i.updateStatus(ctx, func(s *Status) {
				s.LastHash = hash
				s.CountProcessed++
				s.AddError(err)
			})
		}
	}
}

func (i *Importer) WatchUpdates(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer func() { _ = watcher.Close() }()

	for alias, dir := range i.directories {
		if alias == i.directoriesUploadsAlias {
			continue
		}
		if err = watcher.Add(dir); err != nil {
			return fmt.Errorf("add watcher for %q: %w", alias, err)
		}
		log.Printf("starting watcher for %q", dir)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}
			if strings.HasSuffix(event.Name, ".tmp") {
				continue
			}
			dir := filepath.Dir(event.Name)
			dirAlias, ok := i.directories.AliasByPath(dir)
			if !ok {
				continue
			}
			fileName := filepath.Base(event.Name)
			if err := i.EnqueueFile(ctx, dirAlias, dir, fileName, time.Now()); err != nil {
				log.Printf("error enqueueing watcher event %v: %v", event, err)
			}
		}
	}
}

func (i *Importer) importMediaFromFile(ctx context.Context, dirAlias, dir, fileName string, modTime time.Time) (string, error) {
	_, err := i.db.GetDirInfo(ctx, dirAlias, fileName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("getting dir info: %w", err)
	}
	if err == nil {
		return "", nil
	}

	log.Printf("importing new item. alias %q, filename %q", dirAlias, fileName)

	filePath := filepath.Join(dir, fileName)
	raw, err := os.ReadFile(filePath) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}

	media, err := imagery.NewMedia(raw)
	if err != nil {
		return "", fmt.Errorf("decode and hash media: %w", err)
	}

	timestamp := GuessFileCreated(fileName, modTime)

	id, isOld, err := i.insertMedia(ctx, media, timestamp)
	if err != nil {
		return "", fmt.Errorf("insert media: %w", err)
	}
	if err := i.insertDirInfo(ctx, id, dirAlias, fileName); err != nil {
		return "", fmt.Errorf("insert dir info: %w", err)
	}
	for _, f := range i.notifyMediaFuncs {
		f(ctx, media.Hash())
	}

	if isOld {
		return media.Hash(), nil
	}

	if err := i.insertThumbnail(ctx, id, media.Image()); err != nil {
		return "", fmt.Errorf("insert thumbnail: %w", err)
	}
	if err := i.insertBlocks(ctx, id, media.Image()); err != nil {
		return "", fmt.Errorf("insert blocks: %w", err)
	}
	if err := i.db.SetMediaProcessed(ctx, id); err != nil {
		return "", fmt.Errorf("set media processed: %w", err)
	}
	for _, f := range i.notifyMediaFuncs {
		f(ctx, media.Hash())
	}

	return media.Hash(), nil
}

func (i *Importer) scanDirectories(ctx context.Context) error {
	i.updateStatus(ctx, func(s *Status) {
		s.Running = true
		s.CountTotal = 0
		s.CountProcessed = 0
		s.LastHash = ""
		s.Errors = Errors{}
	})
	defer i.updateStatus(ctx, func(s *Status) {
		s.Running = false
	})

	var mediaFiles []*mediaFile
	for alias, dir := range i.directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("listing dir %q: %v", dir, err)
			i.updateStatus(ctx, func(s *Status) {
				s.AddError(fmt.Errorf("listing dir %q: %w", dir, err))
			})
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			info, err := file.Info()
			if err != nil {
				log.Printf("get file info %q: %v", fileName, err)
				i.updateStatus(ctx, func(s *Status) {
					s.AddError(fmt.Errorf("get file info %q: %w", fileName, err))
				})
				continue
			}
			modTime := info.ModTime()
			mediaFiles = append(mediaFiles, &mediaFile{alias, dir, fileName, modTime})
		}
	}

	i.updateStatus(ctx, func(s *Status) {
		s.CountTotal = len(mediaFiles)
	})

	for _, mediaFile := range mediaFiles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case i.jobs <- mediaFile:
		}
	}
	return nil
}

func (i *Importer) updateStatus(ctx context.Context, f func(*Status)) {
	i.statusMu.Lock()
	f(&i.status)
	i.statusMu.Unlock()
	for _, f := range i.notifyProgressFuncs {
		f(ctx)
	}
}

func (i *Importer) insertMedia(ctx context.Context, media imagery.Media, timestamp time.Time) (db.MediaID, bool, error) {
	old, err := i.db.GetMediaByHash(ctx, media.Hash())
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("getting media by hash: %w", err)
	}
	if err == nil {
		return old.ID, true, nil
	}

	_, propDominantColour := imagery.DominantColour(media.Image())

	propBlurhash, err := imagery.CalculateBlurhash(media.Image())
	if err != nil {
		return 0, false, fmt.Errorf("calculate blurhash: %w", err)
	}

	propDimensions := media.Image().Bounds().Size()
	created, err := i.db.CreateMedia(ctx, &db.Media{
		Hash:           media.Hash(),
		Type:           db.MediaType(media.Type()),
		MIME:           media.MIME(),
		Timestamp:      timestamp,
		DimWidth:       propDimensions.X,
		DimHeight:      propDimensions.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// another worker inserted this hash concurrently; treat as already imported
		existing, err := i.db.GetMediaByHash(ctx, media.Hash())
		if err != nil {
			return 0, false, fmt.Errorf("getting raced media: %w", err)
		}
		return existing.ID, true, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("inserting media: %w", err)
	}

	return created.ID, false, nil
}

func (i *Importer) insertBlocks(ctx context.Context, id db.MediaID, image image.Image) error {
	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.ResizeFactor(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := i.defaultEncoder(imageEncoded, imageBig); err != nil {
		return fmt.Errorf("encode scaled and greyed image: %w", err)
	}
	rawBlocks, err := imagery.ExtractText(imageEncoded.Bytes())
	if err != nil {
		return fmt.Errorf("extract image text: %w", err)
	}

	blocks := make([]*db.Block, 0, len(rawBlocks))
	for idx, rawBlock := range rawBlocks {
		if strings.TrimSpace(rawBlock.Word) == "" {
			continue
		}

		rect := imagery.ScaleDownRect(rawBlock.Box)
		blocks = append(blocks, &db.Block{
			MediaID: id,
			Index:   idx,
			MinX:    rect.Min.X,
			MinY:    rect.Min.Y,
			MaxX:    rect.Max.X,
			MaxY:    rect.Max.Y,
			Body:    rawBlock.Word,
		})
	}

	if err := i.db.CreateBlocks(ctx, blocks); err != nil {
		return fmt.Errorf("inserting blocks: %w", err)
	}
	return nil
}

func (i *Importer) insertThumbnail(ctx context.Context, id db.MediaID, image image.Image) error {
	resized := imagery.Resize(image, i.thumbnailWidth, 0)
	dimensions := resized.Bounds().Size()

	var data bytes.Buffer
	if err := i.defaultEncoder(&data, resized); err != nil {
		return fmt.Errorf("encoding thumbnail: %w", err)
	}

	thumbnail := &db.Thumbnail{
		MediaID:   id,
		MIME:      i.defaultMIME,
		DimWidth:  dimensions.X,
		DimHeight: dimensions.Y,
		Timestamp: time.Now(),
		Data:      data.Bytes(),
	}
	if _, err := i.db.CreateThumbnail(ctx, thumbnail); err != nil {
		return fmt.Errorf("insert thumbnail: %w", err)
	}
	return nil
}

func (i *Importer) insertDirInfo(ctx context.Context, id db.MediaID, dirAlias string, fileName string) error {
	dirInfo := &db.DirInfo{
		Filename:       fileName,
		DirectoryAlias: dirAlias,
		MediaID:        id,
	}
	if err := i.db.CreateDirInfo(ctx, dirInfo); err != nil {
		return fmt.Errorf("insert info dir infos: %w", err)
	}
	return nil
}

var fileStampExpr = regexp.MustCompile(`(?:\D|^)(?P<ymd>(?:19|20|21)\d{6})\D?(?P<hms>\d{6})(?:\D|$)`)

func GuessFileCreated(fileName string, modTime time.Time) time.Time {
	fileName = filepath.Base(fileName)
	fileName = strings.TrimPrefix(fileName, "IMG_")
	fileName = strings.TrimPrefix(fileName, "VID_")
	fileName = strings.TrimPrefix(fileName, "img_")
	fileName = strings.TrimPrefix(fileName, "vid_")
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// first try RFC3339
	if guessed, err := time.Parse(time.RFC3339, fileName); err == nil {
		return guessed
	}

	// if that doesn't work, try the date parse library
	if guessed, err := dateparse.ParseLocal(fileName); err == nil {
		return guessed
	}

	// maybe a YYYYMMDD-HHMMSS pattern
	if m := fileStampExpr.FindStringSubmatch(fileName); len(m) > 0 {
		ymd := m[fileStampExpr.SubexpIndex("ymd")]
		hms := m[fileStampExpr.SubexpIndex("hms")]
		guessed, _ := time.Parse("20060102150405", ymd+hms)
		return guessed
	}

	// otherwise, fallback to the file's mod time
	return modTime
}

type mediaFile struct {
	dirAlias string
	dir      string
	fileName string
	modTime  time.Time
}

type Errors []StatusError
type StatusError struct {
	Time  time.Time
	Error error
}

type Status struct {
	Running        bool
	CountTotal     int
	CountProcessed int
	LastHash       string
	Errors         Errors
}

func (s *Status) AddError(err error) {
	if err == nil {
		return
	}
	s.Errors = append(s.Errors, StatusError{
		Time:  time.Now(),
		Error: err,
	})
	if len(s.Errors) > 20 {
		copy(s.Errors, s.Errors[1:])
		s.Errors = s.Errors[:len(s.Errors)-1]
	}
}
