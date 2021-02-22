package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	_ "github.com/lib/pq"

	"go.senan.xyz/socr/backend/controller"
	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/backend/sql"
	"go.senan.xyz/socr/frontend"
)

func main() {
	confListenAddr := mustEnv("SOCR_LISTEN_ADDR")
	confDB := mustEnv("SOCR_DB")
	confHMACSecret := mustEnv("SOCR_HMAC_SECRET")
	confLoginUsername := mustEnv("SOCR_LOGIN_USERNAME")
	confLoginPassword := mustEnv("SOCR_LOGIN_PASSWORD")
	confAPIKey := mustEnv("SOCR_API_KEY")
	confDirs := mustEnvDirs("SOCR_DIR_")

	db, err := db.NewConn(confDB)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}

	importer := &importer.Importer{
		Directories: confDirs,
		DB:          db,
	}

	_ = importer

	ctrl := &controller.Controller{
		Directories: confDirs,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: this?
				return true
			},
		},
		// SocketClientsSettings:   map[*websocket.Conn]struct{}{},
		// SocketClientsScreenshot: map[string]map[*websocket.Conn]struct{}{},
		// SocketUpdatesSettings:   make(chan struct{}),
		// SocketUpdatesScreenshot: make(chan *screenshot.Screenshot),
		HMACSecret:    confHMACSecret,
		LoginUsername: confLoginUsername,
		LoginPassword: confLoginPassword,
		APIKey:        confAPIKey,
		DefaultFormat: imagery.FormatPNG,
	}

	// go ctrl.EmitUpdatesSettings()
	// go ctrl.EmitUpdatesScreenshot()

	r := mux.NewRouter()
	r.Use(ctrl.WithCORS())
	r.Use(ctrl.WithLogging())
	r.HandleFunc("/api/authenticate", ctrl.ServeAuthenticate)
	r.HandleFunc("/api/screenshot/{dir}/{id}/raw", ctrl.ServeScreenshotRaw)
	r.HandleFunc("/api/screenshot/{dir}/{id}", ctrl.ServeScreenshot)
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

func mustEnvDirs(prefix string) map[string]string {
	expr := regexp.MustCompile(`SOCR_DIR_(?P<Alias>[\w_]+)=(?P<Path>.*)`)
	const (
		partFull = iota
		partAlias
		partPath
	)

	dirs := map[string]string{}
	for _, env := range os.Environ() {
		parts := expr.FindStringSubmatch(env)
		if len(parts) != 3 {
			log.Fatalf("invalid screenshot var %s", env)
		}
		dirs[parts[partAlias]] = parts[partPath]
	}
	return dirs
}
