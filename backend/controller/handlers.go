package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

const processedSuffix = ".processed"

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

	screenshot, err := c.ProcessBytes(raw)
	if err != nil {
		http.Error(w, fmt.Sprintf("processing upload: %v", err), 500)
		return
	}

	fmt.Fprintf(w, "%s\n", screenshot.ID)
}

func (c *Controller) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s.png", vars["id"])
	http.ServeFile(w, r, filepath.Join(c.ScreenshotsPath, filename))
}

func (c *Controller) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(c.ImportPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("listing import dir: %v", err), 500)
		return
	}

	go func() {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), processedSuffix) {
				continue
			}

			filePath := filepath.Join(c.ImportPath, file.Name())
			bytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("error reading import: %v", err)
				continue

			}

			screenshot, err := c.ProcessBytes(bytes)
			if err != nil {
				log.Printf("error processing import: %v", err)
				continue

			}

			log.Printf("processed import: %s", screenshot.ID)
		}
	}()
}
