package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

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

	screenshot, err := c.ReadAndIndexBytes(raw)
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
			raw, err := c.ReadRenameImportFile(file)
			if err != nil {
				log.Printf("error processing imported file: %v", err)
				continue
			}

			if raw == nil {
				// file has likely already been imported
				continue
			}

			screenshot, err := c.ReadAndIndexBytes(raw)
			if err != nil {
				log.Printf("processing and indexing: %v", err)
				continue
			}

			log.Printf("processed import: %s", screenshot.ID)

		}
	}()
}

func (c *Controller) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := c.SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket connection: %v", err)
		return
	}

	log.Printf("new socket client: %v", conn.RemoteAddr())
	c.SocketClients[conn] = struct{}{}
}
