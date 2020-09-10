package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"go.senan.xyz/socr/imagery"

	"github.com/blevesearch/bleve"
)

type Controller struct {
	ScreenshotsPath string
	ImportPath      string
	Index           bleve.Index
}

type Screenshot struct {
	ID        string              `json:"id"`
	Timestamp time.Time           `json:"timestamp"`
	Tags      []string            `json:"tags"`
	Processed *imagery.Screenshot `json:"processed"`
}

func (c *Controller) ServeUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	infile, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("read form: %v", err), 500)
		return
	}
	defer infile.Close()

	raw, err := ioutil.ReadAll(infile)
	if err != nil {
		http.Error(w, fmt.Sprintf("read form bytes: %v", err), 500)
		return
	}

	scrotProcessed, err := imagery.ProcessBytes(raw)
	if err != nil {
		http.Error(w, fmt.Sprintf("processing image: %v", err), 500)
		return
	}

	scrotID := IDNew()
	scrotFilename := fmt.Sprintf("%s.%s", scrotID, scrotProcessed.Filetype)
	scrotPath := filepath.Join(c.ScreenshotsPath, scrotFilename)
	if err := ioutil.WriteFile(scrotPath, raw, 0644); err != nil {
		http.Error(w, fmt.Sprintf("write processed bytes: %v", err), 500)
		return
	}

	screenshot := &Screenshot{
		ID:        scrotID,
		Processed: scrotProcessed,
		Timestamp: time.Now(),
		Tags:      []string{},
	}

	if err := c.Index.Index(screenshot.ID, screenshot); err != nil {
		http.Error(w, fmt.Sprintf("indexing screenshot: %v", err), 500)
		return
	}

	fmt.Fprintf(w, "%s\n", scrotFilename)
}

func (c *Controller) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s.png", vars["id"])
	http.ServeFile(w, r, filepath.Join(c.ScreenshotsPath, filename))
}

func (c *Controller) StartImport(w http.ResponseWriter, r *http.Request) {
}
