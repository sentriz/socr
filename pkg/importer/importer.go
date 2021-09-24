package importer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
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

	"go.senan.xyz/socr/pkg/db"
	"go.senan.xyz/socr/pkg/directories"
	"go.senan.xyz/socr/pkg/imagery"
)

type NotifyMediaFunc func(imagery.Hash)
type NotifyProgressFunc func()

type Importer struct {
	db                      *db.DB
	defaultEncoder          imagery.EncodeFunc
	defaultMIME             imagery.MIME
	directories             directories.Directories
	directoriesUploadsAlias string
	thumbnailWidth          uint

	status              Status
	jobs                chan *mediaFile
	notifyMediaFuncs    []NotifyMediaFunc
	notifyProgressFuncs []NotifyProgressFunc
}

func New(
	db *db.DB, defaultEncoder imagery.EncodeFunc, defaultMIME imagery.MIME,
	directories directories.Directories, directoriesUploadsAlias string, thumbnailWidth uint,
) *Importer {
	return &Importer{
		db:                      db,
		defaultEncoder:          defaultEncoder,
		defaultMIME:             defaultMIME,
		directories:             directories,
		directoriesUploadsAlias: directoriesUploadsAlias,
		thumbnailWidth:          thumbnailWidth,

		status: Status{RWMutex: &sync.RWMutex{}},
		jobs:   make(chan *mediaFile),
	}
}

func (i *Importer) AddNotifyMediaFunc(f NotifyMediaFunc) {
	i.notifyMediaFuncs = append(i.notifyMediaFuncs, f)
}

func (i *Importer) AddNotifyProgressFunc(f NotifyProgressFunc) {
	i.notifyProgressFuncs = append(i.notifyProgressFuncs, f)
}

func (i *Importer) ImportMedia(
	fileType imagery.Type, mime imagery.MIME, extension string, image image.Image, hash imagery.Hash,
	dirAlias string, fileName string, timestamp time.Time,
) error {
	id, isOld, err := i.insertMedia(fileType, mime, image, hash, timestamp)
	if err != nil {
		return fmt.Errorf("import media with props: %w", err)
	}
	if err := i.insertDirInfo(id, dirAlias, fileName); err != nil {
		return fmt.Errorf("import dir info: %w", err)
	}
	for _, f := range i.notifyMediaFuncs {
		f(hash)
	}

	if isOld {
		return nil
	}

	if err := i.insertThumbnail(id, image); err != nil {
		return fmt.Errorf("import thumbnail: %w", err)
	}
	if err := i.insertBlocks(id, image); err != nil {
		return fmt.Errorf("import blocks: %w", err)
	}
	if err := i.db.SetMediaProcessed(id); err != nil {
		return fmt.Errorf("set media processed: %w", err)
	}
	for _, f := range i.notifyMediaFuncs {
		f(hash)
	}

	return nil
}

func (i *Importer) ImportMediaFromFile(dirAlias, dir, fileName string, modTime time.Time) (imagery.Hash, error) {
	_, err := i.db.GetDirInfo(dirAlias, fileName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("getting dir info: %w", err)
	}
	if err == nil {
		return "", nil
	}

	log.Printf("importing new item. alias %q, filename %q", dirAlias, fileName)

	filePath := filepath.Join(dir, fileName)
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}

	fileType, mime, extension, image, hash, err := imagery.DecodeAndHash(raw)
	if err != nil {
		return "", fmt.Errorf("decode and hash media: %w", err)
	}

	timestamp := GuessFileCreated(fileName, modTime)

	if err := i.ImportMedia(fileType, mime, extension, image, hash, dirAlias, fileName, timestamp); err != nil {
		return "", fmt.Errorf("importing media: %w", err)
	}

	return hash, nil
}

func (i *Importer) ScanDirectories() error {
	if i.IsRunning() {
		return fmt.Errorf("already running")
	}

	i.updateStatus(func(s *Status) {
		s.Running = true
		s.CountTotal = 0
		s.CountProcessed = 0
		s.LastHash = ""
		s.Errors = Errors{}
	})
	defer i.updateStatus(func(s *Status) {
		s.lastUpdate = time.Time{}
		s.Running = false
	})

	var mediaFiles []*mediaFile
	for alias, dir := range i.directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("listing dir %q: %w", dir, err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			info, err := file.Info()
			if err != nil {
				return fmt.Errorf("get file info %q: %w", fileName, err)
			}
			modTime := info.ModTime()
			mediaFiles = append(mediaFiles, &mediaFile{alias, dir, fileName, modTime})
		}
	}

	i.updateStatus(func(s *Status) {
		s.CountTotal = len(mediaFiles)
	})

	for idx, mediaFile := range mediaFiles {
		i.jobs <- mediaFile
		i.updateStatus(func(s *Status) {
			s.CountProcessed = idx + 1
		})
	}
	return nil
}

func (i *Importer) Status() Status {
	i.status.RLock()
	defer i.status.RUnlock()
	return i.status
}

func (i *Importer) IsRunning() bool {
	i.status.RLock()
	defer i.status.RUnlock()
	return i.status.Running
}

func (i *Importer) StartWorker() {
	for j := range i.jobs {
		hash, err := i.ImportMediaFromFile(j.dirAlias, j.dir, j.fileName, j.modTime)
		i.updateStatus(func(s *Status) {
			s.LastHash = hash
			s.AddError(err)
		})
	}
}

func (i *Importer) updateStatus(f func(*Status)) {
	i.status.Lock()
	defer i.status.Unlock()
	f(&i.status)

	if time.Since(i.status.lastUpdate) > time.Second {
		i.status.lastUpdate = time.Now()
		for _, f := range i.notifyProgressFuncs {
			f()
		}
	}
}

func (i *Importer) insertMedia(
	fileType imagery.Type, mime imagery.MIME, image image.Image, hash imagery.Hash,
	timestamp time.Time,
) (db.MediaID, bool, error) {
	old, err := i.db.GetMediaByHash(string(hash))
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("getting media by hash: %w", err)
	}
	if err == nil {
		return old.ID, true, nil
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, false, fmt.Errorf("calculate blurhash: %w", err)
	}

	propDimensions := image.Bounds().Size()
	new, err := i.db.CreateMedia(&db.Media{
		Hash:           string(hash),
		Type:           db.MediaType(fileType),
		MIME:           string(mime),
		Timestamp:      timestamp,
		DimWidth:       propDimensions.X,
		DimHeight:      propDimensions.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if err != nil {
		return 0, false, fmt.Errorf("inserting media: %w", err)
	}

	return new.ID, false, nil
}

func (i *Importer) insertBlocks(id db.MediaID, image image.Image) error {
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

	if err := i.db.CreateBlocks(blocks); err != nil {
		return fmt.Errorf("inserting blocks: %w", err)
	}
	return nil
}

func (i *Importer) insertThumbnail(id db.MediaID, image image.Image) error {
	resized := imagery.Resize(image, i.thumbnailWidth, 0)
	dimensions := resized.Bounds().Size()

	var data bytes.Buffer
	if err := i.defaultEncoder(&data, resized); err != nil {
		return fmt.Errorf("encoding thumbnail: %w", err)
	}

	thumbnail := &db.Thumbnail{
		MediaID:   id,
		MIME:      string(i.defaultMIME),
		DimWidth:  dimensions.X,
		DimHeight: dimensions.Y,
		Timestamp: time.Now(),
		Data:      data.Bytes(),
	}
	if _, err := i.db.CreateThumbnail(thumbnail); err != nil {
		return fmt.Errorf("insert thumbnail: %w", err)
	}
	return nil
}

func (i *Importer) insertDirInfo(id db.MediaID, dirAlias string, fileName string) error {
	dirInfo := &db.DirInfo{
		Filename:       fileName,
		DirectoryAlias: dirAlias,
		MediaID:        id,
	}
	if _, err := i.db.CreateDirInfo(dirInfo); err != nil {
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

func (i *Importer) WatchUpdates() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	for alias, dir := range i.directories {
		if alias == i.directoriesUploadsAlias {
			continue
		}
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
		dirAlias, ok := i.directories.AliasByPath(dir)
		if !ok {
			continue
		}
		fileName := filepath.Base(event.Name)
		modTime := time.Now()
		if _, err := i.ImportMediaFromFile(dirAlias, dir, fileName, modTime); err != nil {
			log.Printf("error scanning directory item with event %v: %v", event, err)
		}
	}
	return nil
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
	*sync.RWMutex
	lastUpdate     time.Time
	Running        bool
	CountTotal     int
	CountProcessed int
	LastHash       imagery.Hash
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
		s.Errors = s.Errors[1:]
	}
}
