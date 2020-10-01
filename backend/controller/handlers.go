package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.senan.xyz/socr/controller/auth"
	"go.senan.xyz/socr/controller/id"
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

	screenshotID := id.New()
	go func() {
		if _, err := c.ReadAndIndexBytesWithID(raw, screenshotID); err != nil {
			http.Error(w, fmt.Sprintf("processing upload: %v", err), 500)
			return
		}
	}()

	json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{
		ID: screenshotID,
	})
}

func (c *Controller) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	if err := c.IndexImportDirectory(); err != nil {
		http.Error(w, fmt.Sprintf("start import: %v", err), 500)
		return
	}

	json.NewEncoder(w).Encode(struct{}{})
}

func (c *Controller) ServeScreenshotRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s.png", vars["id"])
	http.ServeFile(w, r, filepath.Join(c.ScreenshotsPath, filename))
}

func (c *Controller) ServeScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	query := bleve.NewDocIDQuery([]string{vars["id"]})
	request := bleve.NewSearchRequest(query)
	highlight := bleve.NewHighlight()
	highlight.Fields = []string{
		"blocks.text",
	}
	request.Highlight = highlight
	request.Fields = []string{
		"blocks.text",
		"blocks.position",
		"dimensions.height",
		"dimensions.width",
	}

	resp, err := c.Index.Search(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("searching index: %v", err), 500)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	token := params.Get("token")
	tokenValid := auth.TokenParse(c.HMACSecret, token) == nil

	conn, err := c.SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket connection: %v", err)
		return
	}

	if w := params.Get("want_screenshot_id"); w != "" {
		if _, ok := c.SocketClientsScreenshot[w]; !ok {
			c.SocketClientsScreenshot[w] = map[*websocket.Conn]struct{}{}
		}
		c.SocketClientsScreenshot[w][conn] = struct{}{}
	}

	if w := params.Get("want_settings"); tokenValid && w != "" {
		c.SocketClientsSettings[conn] = struct{}{}
	}

	log.Printf("new socket client: %v", conn.RemoteAddr())
	log.Printf("` clients settings: %v", c.SocketClientsSettings)
	log.Printf("` clients screenshot: %v", c.SocketClientsScreenshot)
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
