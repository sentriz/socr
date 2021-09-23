package server

import (
	"bytes"
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

	"go.senan.xyz/socr"
	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/directories"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/backend/server/auth"
	"go.senan.xyz/socr/backend/server/resp"
)

type Server struct {
	db                      *db.DB
	directories             directories.Directories
	directoriesUploadsAlias string
	socketUpgrader          websocket.Upgrader
	importer                *importer.Importer
	socketClientsScanner    map[*websocket.Conn]struct{}
	socketClientsImporter   map[imagery.Hash]map[*websocket.Conn]struct{}
	hmacSecret              string
	loginUsername           string
	loginPassword           string
	apiKey                  string
	socketMedias            chan imagery.Hash
	socketScannerUpdates    chan struct{}
}

func New(db *db.DB, importr *importer.Importer, directories directories.Directories, uploadsAlias string, hmacSecret, loginUsername, loginPassword, apkKey string) *Server {
	servr := &Server{
		db:                      db,
		directories:             directories,
		directoriesUploadsAlias: uploadsAlias,
		socketUpgrader:          websocket.Upgrader{CheckOrigin: CheckOrigin},
		importer:                importr,
		socketClientsScanner:    map[*websocket.Conn]struct{}{},
		socketClientsImporter:   map[imagery.Hash]map[*websocket.Conn]struct{}{},
		hmacSecret:              hmacSecret,
		loginUsername:           loginUsername,
		loginPassword:           loginPassword,
		apiKey:                  apkKey,
		socketMedias:            make(chan imagery.Hash),
		socketScannerUpdates:    make(chan struct{}),
	}
	importr.AddNotifyMediaFunc(func(hash imagery.Hash) {
		servr.socketMedias <- hash
	})
	importr.AddNotifyProgressFunc(func() {
		servr.socketScannerUpdates <- struct{}{}
	})
	return servr
}

func (s *Server) Router() *mux.Router {
	// begin normal routes
	r := mux.NewRouter()
	r.Use(s.WithCORS())
	r.Use(s.WithLogging())
	r.HandleFunc("/api/authenticate", s.serveAuthenticate)
	r.HandleFunc("/api/media/{hash}/raw", s.serveMediaRaw)
	r.HandleFunc("/api/media/{hash}/thumb", s.serveMediaThumb)
	r.HandleFunc("/api/media/{hash}", s.serveMedia)
	r.HandleFunc("/api/websocket", s.serveWebSocket)

	// begin authenticated routes
	rJWT := r.NewRoute().Subrouter()
	rJWT.Use(s.WithJWT())
	rJWT.HandleFunc("/api/ping", s.servePing)
	rJWT.HandleFunc("/api/start_import", s.serveStartImport)
	rJWT.HandleFunc("/api/about", s.serveAbout)
	rJWT.HandleFunc("/api/directories", s.serveDirectories)
	rJWT.HandleFunc("/api/import_status", s.serveImportStatus)
	rJWT.HandleFunc("/api/search", s.serveSearch)

	// begin api key routes
	rAPIKey := r.NewRoute().Subrouter()
	rAPIKey.Use(s.WithJWTOrAPIKey())
	rAPIKey.HandleFunc("/api/upload", s.serveUpload)

	// frontend fallback route
	r.NotFoundHandler = s.serveFrontend()

	return r
}

func (s *Server) SocketNotifyScannerUpdate() {
	for range s.socketScannerUpdates {
		for client := range s.socketClientsScanner {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(s.socketClientsScanner, client)
				continue
			}
		}
	}
}

func (s *Server) SocketNotifyMedia() {
	for hash := range s.socketMedias {
		for client := range s.socketClientsImporter[hash] {
			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(s.socketClientsImporter[hash], client)
				continue
			}
		}
	}
}

func (s *Server) serveFrontend() http.Handler {
	fs := http.FS(socr.Dist)
	srv := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); err != nil {
			r.URL.Path = "/"
		}
		srv.ServeHTTP(w, r)
	})
}

func (s *Server) servePing(w http.ResponseWriter, r *http.Request) {
	resp.Write(w, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (s *Server) serveUpload(w http.ResponseWriter, r *http.Request) {
	infile, _, err := r.FormFile("i")
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "get form file: %v", err)
		return
	}
	defer infile.Close()
	raw, err := io.ReadAll(infile)
	if err != nil {
		resp.Errorf(w, http.StatusInternalServerError, "read form file: %v", err)
		return
	}
	fileType, mime, extension, image, hash, err := imagery.DecodeAndHash(raw)
	if err != nil {
		resp.Errorf(w, http.StatusInternalServerError, "decoding media: %v", err)
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	uploadsDir := s.directories[s.directoriesUploadsAlias]
	fileName := fmt.Sprintf("%s.%s", timestamp, extension)
	filePath := filepath.Join(uploadsDir, fileName)
	if err := os.WriteFile(filePath, raw, 0600); err != nil {
		resp.Errorf(w, 500, "write upload to disk: %v", err)
		return
	}

	go func() {
		timestamp := time.Now()
		if err := s.importer.ImportMedia(fileType, mime, extension, image, hash, s.directoriesUploadsAlias, fileName, timestamp); err != nil {
			log.Printf("error processing media %s: %v", hash, err)
			return
		}
	}()

	resp.Write(w, struct {
		ID string `json:"id"`
	}{
		ID: string(hash),
	})
}

func (s *Server) serveStartImport(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := s.importer.ScanDirectories(); err != nil {
			log.Printf("error importing: %v", err)
		}
	}()
	resp.Write(w, struct{}{})
}

func (s *Server) serveAbout(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"version":        socr.Version,
		"api key":        s.apiKey,
		"socket clients": len(s.socketClientsScanner),
	}
	for alias, path := range s.directories {
		key := fmt.Sprintf("directory %q", alias)
		settings[key] = path
	}
	resp.Write(w, settings)
}

type DirectoryCount struct {
	*db.DirectoryCount
	IsUploads bool `json:"is_uploads,omitempty"`
}

func (s *Server) serveDirectories(w http.ResponseWriter, r *http.Request) {
	rawCounts, err := s.db.CountDirectories()
	if err != nil {
		resp.Errorf(w, 500, "counting directories by alias: %v", err)
		return
	}

	counts := make([]*DirectoryCount, 0, len(rawCounts))
	for _, raw := range rawCounts {
		counts = append(counts, &DirectoryCount{
			DirectoryCount: raw,
			IsUploads:      raw.DirectoryAlias == s.directoriesUploadsAlias,
		})
	}
	resp.Write(w, counts)
}

func (s *Server) serveMediaRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	if hash == "" {
		resp.Errorf(w, http.StatusBadRequest, "no media hash provided")
		return
	}
	row, err := s.db.GetDirInfoByMediaHash(hash)
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "requested media not found: %v", err)
		return
	}
	directory, ok := s.directories[row.DirectoryAlias]
	if !ok {
		resp.Errorf(w, 500, "media has invalid alias %q", row.DirectoryAlias)
		return
	}
	http.ServeFile(w, r, filepath.Join(directory, row.Filename))
}

func (s *Server) serveMediaThumb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	if hash == "" {
		resp.Errorf(w, http.StatusBadRequest, "no media hash provided")
		return
	}
	row, err := s.db.GetThumbnailByMediaHash(hash)
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "requested media not found: %v", err)
		return
	}
	http.ServeContent(w, r, hash, row.Timestamp, bytes.NewReader(row.Data))
}

func (s *Server) serveMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	if hash == "" {
		resp.Errorf(w, http.StatusBadRequest, "no media hash provided")
		return
	}
	media, err := s.db.GetMediaByHashWithRelations(hash)
	if err != nil {
		resp.Errorf(w, http.StatusBadRequest, "requested media not found: %v", err)
		return
	}
	resp.Write(w, media)
}

type ServeSearchPayload struct {
	Body      string `json:"body"`
	Directory string `json:"directory"`
	Media     string `json:"media"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Sort      struct {
		Field string `json:"field"`
		Order string `json:"order"`
	} `json:"sort"`
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

func (s *Server) serveSearch(w http.ResponseWriter, r *http.Request) {
	var payload ServeSearchPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.Errorf(w, 400, "decode payload: %v", err)
	}
	defer r.Body.Close()

	start := time.Now()
	medias, err := s.db.SearchMedias(db.SearchMediasOptions{
		Body:      payload.Body,
		Offset:    payload.Offset,
		Limit:     payload.Limit,
		SortField: payload.Sort.Field,
		SortOrder: payload.Sort.Order,
		Directory: payload.Directory,
		Media:     db.MediaType(payload.Media),
		DateFrom:  payload.DateFrom,
		DateTo:    payload.DateTo,
	})
	if err != nil {
		resp.Errorf(w, 500, "searching medias: %v", err)
		return
	}

	resp.Write(w, struct {
		Medias []*db.Media   `json:"medias"`
		Took   time.Duration `json:"took"`
	}{
		Medias: medias,
		Took:   time.Since(start),
	})
}

func (s *Server) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	token := params.Get("token")

	conn, err := s.socketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket connection: %v", err)
		return
	}

	if w := params.Get("want_settings"); w != "" {
		if err := auth.TokenParse(s.hmacSecret, token); err != nil {
			return
		}
		s.socketClientsScanner[conn] = struct{}{}
	}
	if w := imagery.Hash(params.Get("want_media_hash")); w != "" {
		if _, ok := s.socketClientsImporter[w]; !ok {
			s.socketClientsImporter[w] = map[*websocket.Conn]struct{}{}
		}
		s.socketClientsImporter[w][conn] = struct{}{}
	}
}

func (s *Server) serveAuthenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.Errorf(w, http.StatusBadRequest, "decode payload: %v", err)
	}

	hasUsername := (payload.Username == s.loginUsername)
	hasPassword := (payload.Password == s.loginPassword)
	if !(hasUsername && hasPassword) {
		resp.Errorf(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	token, err := auth.TokenNew(s.hmacSecret)
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

func (s *Server) serveImportStatus(w http.ResponseWriter, r *http.Request) {
	type respStatusError struct {
		Time  time.Time `json:"time"`
		Error string    `json:"error"`
	}
	type respStatus struct {
		Running        bool               `json:"running"`
		Errors         []*respStatusError `json:"errors"`
		LastHash       string             `json:"last_hash"`
		CountTotal     int                `json:"count_total"`
		CountProcessed int                `json:"count_processed"`
	}

	status := s.importer.Status()
	statusResp := &respStatus{
		Running:        status.Running,
		CountTotal:     status.CountTotal,
		CountProcessed: status.CountProcessed,
		LastHash:       string(status.LastHash),
	}
	for _, err := range status.Errors {
		statusResp.Errors = append(statusResp.Errors, &respStatusError{
			Time:  err.Time,
			Error: err.Error.Error(),
		})
	}

	resp.Write(w, statusResp)
}

// used for socket upgrader
// not checking origin here because currently to become a socket client,
// you must know the hash of the media, or else provide a token for sensitive info.
// if there is a problem with this please let me know
func CheckOrigin(r *http.Request) bool {
	return true
}
