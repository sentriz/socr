package server

import (
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

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/directories"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/backend/scanner"
	"go.senan.xyz/socr/backend/server/auth"
	"go.senan.xyz/socr/backend/server/resp"
	"go.senan.xyz/socr/frontend"
)

type Server struct {
	DB                      *db.DB
	Directories             directories.Directories
	DirectoriesUploadsAlias string
	SocketUpgrader          websocket.Upgrader
	Importer                *importer.Importer
	Scanner                 *scanner.Scanner
	SocketClientsScanner    map[*websocket.Conn]struct{}
	SocketClientsImporter   map[string]map[*websocket.Conn]struct{}
	HMACSecret              string
	LoginUsername           string
	LoginPassword           string
	APIKey                  string
	DefaultFormat           imagery.Format
}

func (c *Server) Router() *mux.Router {
	// begin normal routes
	r := mux.NewRouter()
	r.Use(c.WithCORS())
	r.Use(c.WithLogging())
	r.HandleFunc("/api/authenticate", c.ServeAuthenticate)
	r.HandleFunc("/api/screenshot/{hash}/raw", c.ServeScreenshotRaw)
	r.HandleFunc("/api/screenshot/{hash}", c.ServeScreenshot)
	r.HandleFunc("/api/websocket", c.ServeWebSocket)

	// begin authenticated routes
	rJWT := r.NewRoute().Subrouter()
	rJWT.Use(c.WithJWT())
	rJWT.HandleFunc("/api/ping", c.ServePing)
	rJWT.HandleFunc("/api/start_import", c.ServeStartImport)
	rJWT.HandleFunc("/api/about", c.ServeAbout)
	rJWT.HandleFunc("/api/directories", c.ServeDirectories)
	rJWT.HandleFunc("/api/import_status", c.ServeImportStatus)
	rJWT.HandleFunc("/api/search", c.ServeSearch)

	// begin api key routes
	rAPIKey := r.NewRoute().Subrouter()
	rAPIKey.Use(c.WithJWTOrAPIKey())
	rAPIKey.HandleFunc("/api/upload", c.ServeUpload)

	// frontend fallback route
	frontendFS := http.FS(frontend.FS)
	r.NotFoundHandler = http.FileServer(frontendFS)

	return r
}

func (c *Server) EmitUpdatesScanner() {
	for range c.Scanner.Updates {
		for client := range c.SocketClientsScanner {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClientsScanner, client)
				continue
			}
		}
	}
}

func (c *Server) EmitUpdatesImporter() {
	for id := range c.Importer.Updates {
		for client := range c.SocketClientsImporter[id] {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClientsImporter[id], client)
				continue
			}
		}
	}
}

func (c *Server) ServePing(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (c *Server) ServeUpload(w http.ResponseWriter, r *http.Request) {
	infile, _, err := r.FormFile("i")
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "get form file: %v", err)
		return
	}
	defer infile.Close()
	raw, err := io.ReadAll(infile)
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "read form file: %v", err)
		return
	}
	decoded, err := importer.DecodeImage(raw)
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "decoding screenshot: %v", err)
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	uploadsDir := c.Directories[c.DirectoriesUploadsAlias]
	fileName := fmt.Sprintf("%s.%s", timestamp, decoded.Format.Filetype)
	filePath := filepath.Join(uploadsDir, fileName)
	if err := os.WriteFile(filePath, raw, 0600); err != nil {
		resp.Errorf(w, 500, "write upload to disk: %v", err)
		return
	}

	go func() {
		timestamp := time.Now()
		if err := c.Importer.ImportScreenshot(decoded, timestamp, c.DirectoriesUploadsAlias, fileName); err != nil {
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

func (c *Server) ServeStartImport(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := c.Scanner.ScanDirectories(); err != nil {
			log.Printf("error importing: %v", err)
		}
	}()
	resp.Write(w, struct{}{})
}

func (c *Server) ServeAbout(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"version":        "development",
		"api key":        c.APIKey,
		"socket clients": len(c.SocketClientsScanner),
	}
	for alias, path := range c.Directories {
		key := fmt.Sprintf("directory %q", alias)
		settings[key] = path
	}
	resp.Write(w, settings)
}

type DirectoryCount struct {
	*db.DirectoryCount
	IsUploads bool `json:"is_uploads,omitempty"`
}

func (c *Server) ServeDirectories(w http.ResponseWriter, r *http.Request) {
	rawCounts, err := c.DB.CountDirectories()
	if err != nil {
		resp.Errorf(w, 500, "counting directories by alias: %v", err)
		return
	}

	counts := make([]*DirectoryCount, 0, len(rawCounts))
	for _, raw := range rawCounts {
		counts = append(counts, &DirectoryCount{
			DirectoryCount: raw,
			IsUploads:      raw.DirectoryAlias == c.DirectoriesUploadsAlias,
		})
	}
	resp.Write(w, counts)
}

func (c *Server) ServeScreenshotRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	row, err := c.DB.GetDirInfoByScreenshotHash(vars["hash"])
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "requested screenshot not found: %v", err)
		return
	}
	directory, ok := c.Directories[row.DirectoryAlias]
	if !ok {
		resp.Errorf(w, 500, "screenshot has invalid alias %q", row.DirectoryAlias)
		return
	}
	http.ServeFile(w, r, filepath.Join(directory, row.Filename))
}

func (c *Server) ServeScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	screenshot, err := c.DB.GetScreenshotByHashWithRelations(vars["hash"])
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "requested screenshot not found: %v", err)
		return
	}
	resp.Write(w, screenshot)
}

type ServeSearchPayload struct {
	Body      string `json:"body"`
	Directory string `json:"directory"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Sort      struct {
		Field string `json:"field"`
		Order string `json:"order"`
	} `json:"sort"`
}

func (c *Server) ServeSearch(w http.ResponseWriter, r *http.Request) {
	var payload ServeSearchPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.Errorf(w, 400, "decode payload: %v", err)
	}
	defer r.Body.Close()

	start := time.Now()
	screenshots, err := c.DB.SearchScreenshots(db.SearchScreenshotsOptions{
		Body:      payload.Body,
		Offset:    payload.Offset,
		Limit:     payload.Limit,
		SortField: payload.Sort.Field,
		SortOrder: payload.Sort.Order,
		Directory: payload.Directory,
	})
	if err != nil {
		resp.Errorf(w, 500, "searching screenshots: %v", err)
		return
	}

	resp.Write(w, struct {
		Screenshots []*db.Screenshot `json:"screenshots"`
		Took        time.Duration    `json:"took"`
	}{
		Screenshots: screenshots,
		Took:        time.Since(start),
	})
}

func (c *Server) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
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
		c.SocketClientsScanner[conn] = struct{}{}
	}
	if w := params.Get("want_screenshot_hash"); w != "" {
		if _, ok := c.SocketClientsImporter[w]; !ok {
			c.SocketClientsImporter[w] = map[*websocket.Conn]struct{}{}
		}
		c.SocketClientsImporter[w][conn] = struct{}{}
	}
}

func (c *Server) ServeAuthenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.Errorf(w, http.StatusBadRequest, "decode payload: %v", err)
	}

	hasUsername := (payload.Username == c.LoginUsername)
	hasPassword := (payload.Password == c.LoginPassword)
	if !(hasUsername && hasPassword) {
		resp.Errorf(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	token, err := auth.TokenNew(c.HMACSecret)
	if err != nil {
		resp.Errorf(w, 500, "generating token")
		return
	}

	resp.Write(w, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (c *Server) ServeImportStatus(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		scanner.Status
		Running bool `json:"running"`
	}{
		Status:  c.Scanner.Status,
		Running: c.Scanner.IsRunning(),
	})
}
