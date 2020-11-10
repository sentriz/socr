package controller

import (
	"bytes"
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

	"go.senan.xyz/socr/controller/id"
	"go.senan.xyz/socr/imagery"

	"github.com/araddon/dateparse"
	"github.com/blevesearch/bleve"
	"github.com/gorilla/websocket"
)

type Screenshot struct {
	ID             string           `json:"id"`
	Timestamp      time.Time        `json:"timestamp"`
	Filetype       imagery.Filetype `json:"filetype"`
	Tags           []ScreenshotTag  `json:"tags"`
	Dimensions     Dimensions       `json:"dimensions"`
	Blocks         []*Block         `json:"blocks"`
	DominantColour string           `json:"dominant_colour"`
	Blurhash       string           `json:"blurhash"`
}

type ScreenshotTag string

const (
	ScreenshotTagImported ScreenshotTag = "imported"
	ScreenshotTagDesktop  ScreenshotTag = "desktop"
	ScreenshotTagMobile   ScreenshotTag = "mobile"
)

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Block struct {
	// [x1 y1 x2 y2]
	Position [4]int `json:"position"`
	Text     string `json:"text"`
}

type ImportStatusError struct {
	Time  time.Time `json:"time"`
	Error string    `json:"error"`
}

type ImportStatus struct {
	Errors         []ImportStatusError `json:"errors,omitempty"`
	LastID         string              `json:"last_id,omitempty"`
	CountProcessed int                 `json:"count_processed"`
	CountTotal     int                 `json:"count_total"`
}

func (s *ImportStatus) AddError(err error) {
	s.Errors = append(s.Errors, ImportStatusError{
		Time:  time.Now(),
		Error: err.Error(),
	})
}

type Controller struct {
	ScreenshotsPath         string
	ImportPath              string
	ImportRunning           *int32
	ImportStatus            ImportStatus
	Index                   bleve.Index
	SocketUpgrader          websocket.Upgrader
	SocketClientsSettings   map[*websocket.Conn]struct{}
	SocketClientsScreenshot map[string]map[*websocket.Conn]struct{}
	SocketUpdatesSettings   chan struct{}
	SocketUpdatesScreenshot chan *Screenshot
	HMACSecret              string
	LoginUsername           string
	LoginPassword           string
	APIKey                  string
	DefaultFormat           imagery.Format
}

func (c *Controller) ImportIsRunning() bool { return atomic.LoadInt32(c.ImportRunning) == 1 }
func (c *Controller) ImportSetRunning()     { atomic.StoreInt32(c.ImportRunning, 1) }
func (c *Controller) ImportSetFinished()    { atomic.StoreInt32(c.ImportRunning, 0) }

func (c *Controller) ReadAndIndexBytes(raw []byte) (*Screenshot, error) {
	return c.ReadAndIndexBytesWithIDTime(raw, id.New(), time.Now())
}

func (c *Controller) ReadAndIndexBytesWithID(raw []byte, scrotID string) (*Screenshot, error) {
	return c.ReadAndIndexBytesWithIDTime(raw, scrotID, time.Now())
}

func (c *Controller) ReadAndIndexBytesWithTime(raw []byte, timestamp time.Time) (*Screenshot, error) {
	return c.ReadAndIndexBytesWithIDTime(raw, id.New(), timestamp)
}

func (c *Controller) ReadAndIndexBytesWithIDTime(raw []byte, scrotID string, timestamp time.Time) (*Screenshot, error) {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unrecognised format: %s", mime)
	}

	scrotPath := filepath.Join(c.ScreenshotsPath, scrotID)
	if err := ioutil.WriteFile(scrotPath, raw, 0644); err != nil {
		return nil, fmt.Errorf("write processed bytes: %w", err)
	}

	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return nil, fmt.Errorf("decoding: %s", mime)
	}

	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := c.DefaultFormat.Encode(imageEncoded, imageBig); err != nil {
		return nil, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	scrotBlocksOrig, err := imagery.ExtractText(imageEncoded.Bytes(), imagery.ScaleFactor)
	if err != nil {
		return nil, fmt.Errorf("extract image text: %w", err)
	}

	scrotBlocks := []*Block{}
	for _, block := range scrotBlocksOrig {
		scrotBlocks = append(scrotBlocks, &Block{
			Position: imagery.ScaleDownRect(block.Box),
			Text:     block.Word,
		})
	}

	_, scrotDominantColour := imagery.DominantColour(image)
	scrotBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	screenshot := &Screenshot{
		ID:       scrotID,
		Filetype: format.Filetype,
		Dimensions: Dimensions{
			Width:  image.Bounds().Size().X,
			Height: image.Bounds().Size().Y,
		},
		DominantColour: scrotDominantColour,
		Blurhash:       scrotBlurhash,
		Blocks:         scrotBlocks,
		Timestamp:      timestamp,
		Tags: []ScreenshotTag{
			ScreenshotTagImported,
		},
	}

	if err := c.Index.Index(screenshot.ID, screenshot); err != nil {
		return nil, fmt.Errorf("indexing screenshot: %w", err)
	}

	return screenshot, nil
}

func procSuffixHas(in string) bool   { return strings.HasSuffix(in, ".processed") }
func procSuffixAdd(in string) string { return fmt.Sprintf("%s.processed", in) }

func (c *Controller) IndexImportFile(file os.FileInfo) (*Screenshot, error) {
	filePath := filepath.Join(c.ImportPath, file.Name())
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading from disk: %v", err)
	}

	if err := os.Rename(filePath, procSuffixAdd(filePath)); err != nil {
		return nil, fmt.Errorf("renaming: %v", err)
	}

	screenshotTimestamp := guessFileCreated(file)
	screenshot, err := c.ReadAndIndexBytesWithTime(bytes, screenshotTimestamp)
	if err != nil {
		return nil, fmt.Errorf("processing and indexing: %v", err)
	}

	return screenshot, nil
}

func (c *Controller) IndexImportDirectory() error {
	if c.ImportIsRunning() {
		return fmt.Errorf("already importing")
	}

	files, err := ioutil.ReadDir(c.ImportPath)
	if err != nil {
		return fmt.Errorf("listing import dir: %w", err)
	}

	go c.IndexImportDirectoryProcess(files)

	return nil
}

func (c *Controller) IndexImportDirectoryProcess(files []os.FileInfo) {
	c.ImportSetRunning()
	defer c.ImportSetFinished()

	c.ImportStatus = ImportStatus{}
	c.SocketUpdatesSettings <- struct{}{}

	var nonProcessed []os.FileInfo
	for _, file := range files {
		if !procSuffixHas(file.Name()) {
			nonProcessed = append(nonProcessed, file)
		}
	}

	if len(nonProcessed) == 0 {
		c.ImportStatus.AddError(errors.New("no more file left in import dir"))
		c.SocketUpdatesSettings <- struct{}{}
		return
	}

	for i, file := range nonProcessed {
		screenshot, err := c.IndexImportFile(file)
		if err != nil {
			c.ImportStatus.AddError(err)
			c.SocketUpdatesSettings <- struct{}{}
			continue
		}

		c.ImportStatus.LastID = screenshot.ID
		c.ImportStatus.CountProcessed = i + 1
		c.ImportStatus.CountTotal = len(nonProcessed)
		c.SocketUpdatesSettings <- struct{}{}
	}
}

func (c *Controller) EmitUpdatesSettings() error {
	for range c.SocketUpdatesSettings {
		for client := range c.SocketClientsSettings {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClientsSettings, client)
				continue
			}
		}
	}
	return nil
}

func (c *Controller) EmitUpdatesScreenshot() error {
	for screenshot := range c.SocketUpdatesScreenshot {
		for client := range c.SocketClientsScreenshot[screenshot.ID] {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClientsScreenshot[screenshot.ID], client)
				continue
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
