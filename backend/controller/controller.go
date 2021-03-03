package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"go.senan.xyz/socr/backend/controller/auth"
	"go.senan.xyz/socr/backend/controller/resp"
	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
)

type Controller struct {
	DB                      *db.Conn
	Directories             map[string]string
	SocketUpgrader          websocket.Upgrader
	Importer                *importer.Importer
	SocketClientsSettings   map[*websocket.Conn]struct{}
	SocketClientsScreenshot map[int64]map[*websocket.Conn]struct{}
	HMACSecret              string
	LoginUsername           string
	LoginPassword           string
	APIKey                  string
	DefaultFormat           imagery.Format
}

func (c *Controller) EmitUpdatesSettings() error {
	for range c.Importer.UpdatesScan {
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
	for screenshot := range c.Importer.UpdatesScreenshot {
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

func (c *Controller) ServePing(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, resp.Status{Status: "ok"})
}

func (c *Controller) ServeUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	infile, _, err := r.FormFile("i")
	if err != nil {
		resp.Error(w, 500, "read form: %v", err)
		return
	}
	defer infile.Close()

	raw, err := ioutil.ReadAll(infile)
	if err != nil {
		resp.Error(w, 500, "read form bytes: %v", err)
		return
	}

	hash, err := hasher.Hash(raw)
	if err != nil {
		resp.Error(w, 500, "hash screenshot: %v", err)
		return
	}

	go func() {
		timestamp := time.Now()
		if _, err := c.Importer.ImportScreenshot(hash, timestamp, "", "", raw); err != nil {
			log.Printf("error processing screenshot %d: %v", hash, err)
			return
		}
	}()

	resp.Write(w, resp.ID{ID: hasher.ID(hash)})
}

func (c *Controller) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := c.Importer.ScanDirectories(); err != nil {
			log.Printf("error importing: %v", err)
		}
	}()
}

func (c *Controller) ServeAbout(w http.ResponseWriter, r *http.Request) {
	screenshotsCount, err := c.DB.CountDirectoriesByAlias(context.Background())
	if err != nil {
		resp.Error(w, 500, "counting directories by alias: %v", err)
		return
	}

	resp.Write(w, resp.About{
		Version:          "development",
		APIKey:           c.APIKey,
		SocketClients:    len(c.SocketClientsSettings),
		ScreenshotsCount: screenshotsCount,
	})
}

func (c *Controller) ServeScreenshotRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID, ok := vars["id"]
	if !ok {
		resp.Error(w, http.StatusBadRequest, "please provide an `id` parameter")
		return
	}

	id, err := hasher.Parse(strID)
	if !ok {
		resp.Error(w, http.StatusBadRequest, "couldn't parse provided id. %v", err)
		return
	}

	screenshot, err := c.DB.GetScreenshotByID(context.Background(), id)
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "provided screenshot not found. %v", err)
		return
	}

	directory, ok := c.Directories[screenshot.DirectoryAlias]
	if !ok {
		resp.Error(w, 500, "screenshot has invalid alias %q", screenshot.DirectoryAlias)
		return
	}

	http.ServeFile(w, r, filepath.Join(directory, screenshot.Filename))
}

func (c *Controller) ServeScreenshot(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// query := bleve.NewDocIDQuery([]string{vars["id"]})
	// request := bleve.NewSearchRequest(query)
	// request.Fields = index.BaseSearchFields

	// resp, err := c.Index.Search(request)
	// if err != nil {
	// 	resp.Error(w, fmt.Sprintf("searching index: %v", err), 500)
	// 	return
	// }

	// json.NewEncoder(w).Encode(resp)
}

type ServeSearchPayload struct {
	Size int      `json:"size"`
	From int      `json:"from"`
	Sort []string `json:"sort"`
	Term string   `json:"term"`
}

func (c *Controller) ServeSearch(w http.ResponseWriter, r *http.Request) {
	var payload ServeSearchPayload
	json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()

	screenshots, err := c.DB.SearchScreenshots(context.Background(), db.SearchScreenshotsParams{
		Body: payload.Term,
		Off:  int32(payload.From),
		Lim:  int32(payload.Size),
	})
	if err != nil {
		resp.Error(w, 500, "searching screenshots: %v", err)
		return
	}

	response, err := resp.NewScreenshots(screenshots)
	if err != nil {
		resp.Error(w, 500, "create response: %v", err)
	}

	resp.Write(w, response)
}

func (c *Controller) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	token := params.Get("token")

	conn, err := c.SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket connection: %v", err)
		return
	}

	if w := params.Get("want_settings"); w != "" {
		if err := auth.TokenParse(c.HMACSecret, token); err != nil {
			return
		}
		c.SocketClientsSettings[conn] = struct{}{}
	}

	if w := params.Get("want_screenshot_id"); w != "" {
		id, err := hasher.Parse(w)
		if err != nil {
			return
		}
		if _, ok := c.SocketClientsScreenshot[id]; !ok {
			c.SocketClientsScreenshot[id] = map[*websocket.Conn]struct{}{}
		}
		c.SocketClientsScreenshot[id][conn] = struct{}{}
	}
}

func (c *Controller) ServeAuthenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.Error(w, http.StatusBadRequest, "parse payload: %v", err)
		return
	}

	hasUsername := (payload.Username == c.LoginUsername)
	hasPassword := (payload.Password == c.LoginPassword)
	if !(hasUsername && hasPassword) {
		resp.Error(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	token, err := auth.TokenNew(c.HMACSecret)
	if err != nil {
		resp.Error(w, 500, "generating token")
		return
	}

	resp.Write(w, resp.Token{
		Token: token,
	})
}

func (c *Controller) ServeImportStatus(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, resp.ImportStatus{
		// ImportStatus: c.ImportStatus,
		// Running:      c.ImportIsRunning(),
	})
}
