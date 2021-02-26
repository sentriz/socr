package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	_ "github.com/lib/pq"

	"go.senan.xyz/socr/backend/controller"
	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/frontend"
)

func main() {
	const uploadsKey = "uploads"

	confListenAddr := mustEnv("SOCR_LISTEN_ADDR")
	confDBDSN := mustEnv("SOCR_DB_DSN")
	confHMACSecret := mustEnv("SOCR_HMAC_SECRET")
	confLoginUsername := mustEnv("SOCR_LOGIN_USERNAME")
	confLoginPassword := mustEnv("SOCR_LOGIN_PASSWORD")
	confAPIKey := mustEnv("SOCR_API_KEY")
	confDirs := parseEnvDirs("SOCR_DIR_")

	if _, ok := confDirs["uploads"]; !ok {
		log.Fatalf("please provide an uploads directory")
	}

	for alias, path := range confDirs {
		log.Printf("using directory alias %q path %q", alias, path)
	}

	dbConn, err := db.NewConn(confDBDSN)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	defer dbConn.Close()

	importr := &importer.Importer{
		Running:               new(int32),
		Directories:           confDirs,
		DirectoriesUploadsKey: uploadsKey,
		DB:                    dbConn,
		UpdatesScan:           make(chan struct{}),
		UpdatesScreenshot:     make(chan *db.Screenshot),
	}

	ctrl := &controller.Controller{
		Directories: confDirs,
		DB:          dbConn,
		Importer:    importr,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: this?
				return true
			},
		},
		SocketClientsSettings:   map[*websocket.Conn]struct{}{},
		SocketClientsScreenshot: map[hasher.ID]map[*websocket.Conn]struct{}{},
		HMACSecret:              confHMACSecret,
		LoginUsername:           confLoginUsername,
		LoginPassword:           confLoginPassword,
		APIKey:                  confAPIKey,
		DefaultFormat:           imagery.FormatPNG,
	}

	go ctrl.EmitUpdatesSettings()
	go ctrl.EmitUpdatesScreenshot()

	r := mux.NewRouter()
	r.Use(ctrl.WithCORS())
	r.Use(ctrl.WithLogging())
	r.HandleFunc("/api/authenticate", ctrl.ServeAuthenticate)
	r.HandleFunc("/api/screenshot/{id}/raw", ctrl.ServeScreenshotRaw)
	r.HandleFunc("/api/screenshot/{id}", ctrl.ServeScreenshot)
	r.HandleFunc("/api/websocket", ctrl.ServeWebSocket)

	frontendFS := http.FS(frontend.FS)
	r.NotFoundHandler = http.FileServer(frontendFS)

	// begin authenticated routes
	rJWT := r.NewRoute().Subrouter()
	rJWT.Use(ctrl.WithJWT())
	rJWT.HandleFunc("/api/ping", ctrl.ServePing)
	rJWT.HandleFunc("/api/start_import", ctrl.ServeStartImport)
	rJWT.HandleFunc("/api/about", ctrl.ServeAbout)
	rJWT.HandleFunc("/api/import_status", ctrl.ServeImportStatus)
	rJWT.HandleFunc("/api/search", ctrl.ServeSearch)

	rAPIKey := r.NewRoute().Subrouter()
	rAPIKey.Use(ctrl.WithAPIKey())
	rAPIKey.HandleFunc("/api/upload", ctrl.ServeUpload)

	server := http.Server{
		Addr:    confListenAddr,
		Handler: r,
	}

	log.Printf("listening on %q", confListenAddr)
	log.Fatalf("error starting server: %v", server.ListenAndServe())
}

func mustEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("please provide a %q", key)
	return ""
}

func parseEnvDirs(prefix string) map[string]string {
	expr := regexp.MustCompile(prefix + `(?P<Alias>[\w_]+)=(?P<Path>.*)`)
	const (
		partFull = iota
		partAlias
		partPath
	)
	dirs := map[string]string{}
	for _, env := range os.Environ() {
		parts := expr.FindStringSubmatch(env)
		if len(parts) != 3 {
			continue
		}
		alias := strings.ToLower(parts[partAlias])
		dirs[alias] = parts[partPath]
	}
	return dirs
}
