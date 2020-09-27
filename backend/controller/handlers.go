package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.senan.xyz/socr/controller/auth"
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

	if _, err := c.ReadAndIndexBytes(raw); err != nil {
		http.Error(w, fmt.Sprintf("processing upload: %v", err), 500)
		return
	}

	json.NewEncoder(w).Encode(struct{}{})
}

func (c *Controller) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	if err := c.IndexImportDirectory(); err != nil {
		http.Error(w, fmt.Sprintf("start import: %v", err), 500)
		return
	}

	json.NewEncoder(w).Encode(struct{}{})
}

func (c *Controller) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s.png", vars["id"])
	http.ServeFile(w, r, filepath.Join(c.ScreenshotsPath, filename))
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

func (c *Controller) ServeAuthenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("parse payload: %v", err), 200)
		return
	}

	hasUsername := (payload.Username == c.LoginUsername)
	hasPassword := (payload.Password == c.LoginPassword)

	if !(hasUsername && hasPassword) {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}

	token, err := auth.TokenNew(c.HMACSecret)
	if err != nil {
		http.Error(w, "generating token", 500)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
