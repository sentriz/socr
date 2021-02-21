package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.senan.xyz/socr/controller/id"
	"go.senan.xyz/socr/imagery"

	"github.com/araddon/dateparse"
	"github.com/blevesearch/bleve"
	"github.com/gorilla/websocket"
)

type Screenshot struct {
	ID        string           `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
	Filetype  imagery.Filetype `json:"filetype"`
	Directory string           `json:"directory"`
	Tags      []string         `json:"tags"`
	*imagery.Properties
}

type Controller struct {
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

func (c *Controller) ReadAndIndex(raw []byte) (*Screenshot, error) {
	return c.ReadAndIndexWithIDTime(raw, id.New(), time.Now())
}

func (c *Controller) ReadAndIndexWithID(raw []byte, scrotID string) (*Screenshot, error) {
	return c.ReadAndIndexWithIDTime(raw, scrotID, time.Now())
}

func (c *Controller) ReadAndIndexWithTime(raw []byte, timestamp time.Time) (*Screenshot, error) {
	return c.ReadAndIndexWithIDTime(raw, id.New(), timestamp)
}

func (c *Controller) ReadAndIndexWithIDTime(raw []byte, scrotID string, timestamp time.Time) (*Screenshot, error) {
	properties, err := imagery.Process(raw)
	if err != nil {
		return nil, fmt.Errorf("getting image properties: %w", err)
	}

	screenshot := &Screenshot{
		ID:         scrotID,
		Filetype:   properties.Format.Filetype,
		Timestamp:  timestamp,
		Properties: properties,
	}

	if err := c.Index.Index(screenshot.ID, screenshot); err != nil {
		return nil, fmt.Errorf("indexing screenshot: %w", err)
	}

	return screenshot, nil
}

func (c *Controller) IndexImportFile(file os.FileInfo) (*Screenshot, error) {
	filePath := filepath.Join(c.ImportPath, file.Name())
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading from disk: %v", err)
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
