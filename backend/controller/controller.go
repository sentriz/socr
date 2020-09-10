package controller

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"go.senan.xyz/socr/imagery"

	"github.com/blevesearch/bleve"
)

type Screenshot struct {
	ID        string              `json:"id"`
	Timestamp time.Time           `json:"timestamp"`
	Tags      []string            `json:"tags"`
	Processed *imagery.Screenshot `json:"processed"`
}

type Controller struct {
	ScreenshotsPath string
	ImportPath      string
	Index           bleve.Index
}

func (c *Controller) ProcessBytes(bytes []byte) (*Screenshot, error) {
	scrotProcessed, err := imagery.ProcessBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("processing image: %w", err)
	}

	scrotID := IDNew()
	scrotFilename := fmt.Sprintf("%s.%s", scrotID, scrotProcessed.Filetype)
	scrotPath := filepath.Join(c.ScreenshotsPath, scrotFilename)
	if err := ioutil.WriteFile(scrotPath, bytes, 0644); err != nil {
		return nil, fmt.Errorf("write processed bytes: %w", err)
	}

	screenshot := &Screenshot{
		ID:        scrotID,
		Processed: scrotProcessed,
		Timestamp: time.Now(),
		Tags:      []string{},
	}

	if err := c.Index.Index(screenshot.ID, screenshot); err != nil {
		return nil, fmt.Errorf("indexing screenshot: %w", err)
	}

	return screenshot, nil
}
