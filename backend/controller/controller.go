package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"go.senan.xyz/socr/backend/controller/auth"
	"go.senan.xyz/socr/backend/controller/resp"
	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
)

type Controller struct {
	DB                      *db.DB
	Directories             map[string]string
	DirectoriesUploadsKey   string
	SocketUpgrader          websocket.Upgrader
	Importer                *importer.Importer
	SocketClientsSettings   map[*websocket.Conn]struct{}
	SocketClientsScreenshot map[string]map[*websocket.Conn]struct{}
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
	for id := range c.Importer.UpdatesScreenshot {
		for client := range c.SocketClientsScreenshot[id] {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClientsScreenshot[id], client)
				continue
			}
		}
	}
	return nil
}

func (c *Controller) ServePing(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (c *Controller) ServeUpload(w http.ResponseWriter, r *http.Request) {
	infile, _, err := r.FormFile("i")
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "get form file: %v", err)
		return
	}
	defer infile.Close()
	raw, err := io.ReadAll(infile)
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "read form file: %v", err)
		return
	}
	decoded, err := importer.DecodeImage(raw)
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "decoding screenshot: %v", err)
		return
	}

	uploadsDir := c.Directories[c.DirectoriesUploadsKey]
	timestamp := time.Now().Format(time.RFC3339)
	fileName := fmt.Sprintf("%s.%s", timestamp, decoded.Format.Filetype)
	filePath := filepath.Join(uploadsDir, fileName)
	if err := os.WriteFile(filePath, raw, 0644); err != nil {
		resp.Error(w, 500, "write upload to disk: %v", err)
		return
	}

	go func() {
		timestamp := time.Now()
		if err := c.Importer.ImportScreenshot(decoded, timestamp, c.DirectoriesUploadsKey, fileName); err != nil {
			log.Printf("error processing screenshot %s: %v", decoded.Hash, err)
			return
		}
	}()

	resp.Write(w, struct {
		ID string `json:"id"`
	}{
		ID: decoded.Hash,
	})
}

func (c *Controller) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := c.Importer.ScanDirectories(); err != nil {
			log.Printf("error importing: %v", err)
		}
	}()
	resp.Write(w, struct{}{})
}

func (c *Controller) ServeAbout(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		Version       string `json:"version"`
		APIKey        string `json:"api_key"`
		SocketClients int    `json:"socket_clients"`
	}{
		Version:       "development",
		APIKey:        c.APIKey,
		SocketClients: len(c.SocketClientsSettings),
	})
}

func (c *Controller) ServeDirectories(w http.ResponseWriter, r *http.Request) {
	screenshotsCount, err := c.DB.CountDirectoriesByAlias(context.Background())
	if err != nil {
		resp.Error(w, 500, "counting directories by alias: %v", err)
		return
	}
	resp.Write(w, screenshotsCount)
}

func (c *Controller) ServeScreenshotRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	row, err := c.DB.GetScreenshotPathByHash(context.Background(), vars["hash"])
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "requested screenshot not found: %v", err)
		return
	}
	directory, ok := c.Directories[row.DirectoryAlias]
	if !ok {
		resp.Error(w, 500, "screenshot has invalid alias %q", row.DirectoryAlias)
		return
	}
	http.ServeFile(w, r, filepath.Join(directory, row.Filename))
}

func (c *Controller) ServeScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	screenshot, err := c.DB.GetScreenshotWithBlocksByHash(context.Background(), vars["hash"])
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "requested screenshot not found: %v", err)
		return
	}
	resp.Write(w, screenshot)
}

type ServeSearchPayload struct {
	Size int    `json:"size"`
	From int    `json:"from"`
	Term string `json:"term"`
	Sort struct {
		Field string `json:"field"`
		Order string `json:"order"`
	} `json:"sort"`
}

func (c *Controller) ServeSearch(w http.ResponseWriter, r *http.Request) {
	var payload ServeSearchPayload
	json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()

	start := time.Now()
	var screenshots interface{}
	var err error
	switch {
	case payload.Term != "":
		screenshots, err = c.DB.SearchScreenshots(context.Background(), db.SearchScreenshotsParams{
			Body:   payload.Term,
			Offset: payload.From,
			Limit:  payload.Size,
		})
	default:
		screenshots, err = c.DB.GetAllScreenshots(context.Background(), db.GetAllScreenshotsParams{
			Offset:    payload.From,
			Limit:     payload.Size,
			SortField: payload.Sort.Field,
			SortOrder: payload.Sort.Order,
		})
	}
	if err != nil {
		resp.Error(w, 500, "searching screenshots: %v", err)
		return
	}

	resp.Write(w, struct {
		Screenshots interface{}   `json:"screenshots"`
		Took        time.Duration `json:"took"`
	}{
		Screenshots: screenshots,
		Took:        time.Since(start),
	})
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
	if w := params.Get("want_screenshot_hash"); w != "" {
		if _, ok := c.SocketClientsScreenshot[w]; !ok {
			c.SocketClientsScreenshot[w] = map[*websocket.Conn]struct{}{}
		}
		c.SocketClientsScreenshot[w][conn] = struct{}{}
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

	resp.Write(w, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (c *Controller) ServeImportStatus(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		importer.Status
		Running bool `json:"running"`
	}{
		Status:  c.Importer.Status,
		Running: c.Importer.IsRunning(),
	})
}
