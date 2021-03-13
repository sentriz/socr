package importer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
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

	raw, err := os.ReadFile(item.fileName)
	if err != nil {
		return "", fmt.Errorf("open file: %v", err)
	}

	decoded, err := DecodeImage(raw)
	if err != nil {
		return "", fmt.Errorf("decode screenshot: %v", err)
	}

	timestamp := guessFileCreated(item.fileName, item.modTime)
	if err := i.ImportScreenshot(decoded, timestamp, item.dirAlias, item.fileName); err != nil {
		return "", fmt.Errorf("importing screenshot: %v", err)
	}

	return decoded.Hash, nil
}

func (i *Importer) ImportScreenshot(decoded *Decoded, timestamp time.Time, dirAlias, fileName string) error {
	// insert screenshot and dir info, alert clients with update
	id, isOld, err := i.importScreenshot(decoded.Hash, decoded.Image, timestamp)
	if err != nil {
		return fmt.Errorf("props and blocks: %w", err)
	}
	if err := i.importScreenshotDirInfo(id, dirAlias, fileName); err != nil {
		return fmt.Errorf("dir info: %w", err)
	}
	i.UpdatesScreenshot <- decoded.Hash

	if isOld {
		return nil
	}

	// insert blocks, alert clients with update
	if err := i.importScreenshotBlocks(id, decoded.Image); err != nil {
		return fmt.Errorf("dir info: %w", err)
	}
	i.UpdatesScreenshot <- decoded.Hash

	return nil
}

func (i *Importer) importScreenshot(hash string, image image.Image, timestamp time.Time) (int, bool, error) {
	old, err := i.DB.GetScreenshotByHash(context.Background(), hash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("getting screenshot by hash: %w", err)
	}
	if err == nil {
		return old.ID, true, nil
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, false, fmt.Errorf("calculate blurhash: %w", err)
	}

	propSize := image.Bounds().Size()
	new, err := i.DB.CreateScreenshot(context.Background(), db.CreateScreenshotParams{
		Hash:           hash,
		Timestamp:      timestamp,
		DimWidth:       propSize.X,
		DimHeight:      propSize.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if err != nil {
		return 0, false, fmt.Errorf("inserting screenshot: %w", err)
	}

	return new.ID, false, nil
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
			MinX:         rect.Min.X,
			MinY:         rect.Min.Y,
			MaxX:         rect.Max.X,
			MaxY:         rect.Max.Y,
			Body:         block.Word,
		})
	}
	results := i.DB.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return fmt.Errorf("end transaction: %w", err)
	}

	return nil
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

type Decoded struct {
	Hash   string
	Image  image.Image
	Data   []byte
	Format imagery.Format
}

// DecodeImage takes a raw byte slice of an image (png/jpg/etc) decodes it using an appropriate format
// and encodes it again. Encoding and decoding makes sure the hash will be the same for the same
// image given different sources. clipboard/filesystem/etc
func DecodeImage(raw []byte) (*Decoded, error) {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unknown image mime %s", mime)
	}
	image, err := format.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decoding image %w", err)
	}
	dataBuff := &bytes.Buffer{}
	if err := format.Encode(dataBuff, image); err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	data := dataBuff.Bytes()
	hash := hasher.Hash(data)
	return &Decoded{
		Hash:   hash,
		Image:  image,
		Data:   data,
		Format: format,
	}, err
}
