package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/backend/scanner"
	"go.senan.xyz/socr/backend/server"
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

	if _, ok := confDirs[uploadsKey]; !ok {
		log.Fatalf("please provide an uploads directory")
	}

	for alias, path := range confDirs {
		log.Printf("using directory alias %q path %q", alias, path)
	}

	dbConn, err := db.New(confDBDSN)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	defer dbConn.Close()

	importr := &importer.Importer{
		DB:      dbConn,
		Updates: make(chan string),
	}

	scanr := &scanner.Scanner{
		Running:     new(int32),
		Directories: confDirs,
		DB:          dbConn,
		Importer:    importr,
		Updates:     make(chan struct{}),
	}

	servr := &server.Server{
		Directories:           confDirs,
		DirectoriesUploadsKey: uploadsKey,
		DB:                    dbConn,
		Importer:              importr,
		Scanner:               scanr,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: this?
				return true
			},
		},
		SocketClientsScanner:  map[*websocket.Conn]struct{}{},
		SocketClientsImporter: map[string]map[*websocket.Conn]struct{}{},
		HMACSecret:            confHMACSecret,
		LoginUsername:         confLoginUsername,
		LoginPassword:         confLoginPassword,
		APIKey:                confAPIKey,
		DefaultFormat:         imagery.FormatPNG,
	}

	go servr.EmitUpdatesScanner()
	go servr.EmitUpdatesImporter()

	// begin normal routes
	r := mux.NewRouter()
	r.Use(servr.WithCORS())
	r.Use(servr.WithLogging())
	r.HandleFunc("/api/authenticate", servr.ServeAuthenticate)
	r.HandleFunc("/api/screenshot/{hash}/raw", servr.ServeScreenshotRaw)
	r.HandleFunc("/api/screenshot/{hash}", servr.ServeScreenshot)
	r.HandleFunc("/api/websocket", servr.ServeWebSocket)

	// begin authenticated routes
	rJWT := r.NewRoute().Subrouter()
	rJWT.Use(servr.WithJWT())
	rJWT.HandleFunc("/api/ping", servr.ServePing)
	rJWT.HandleFunc("/api/start_import", servr.ServeStartImport)
	rJWT.HandleFunc("/api/about", servr.ServeAbout)
	rJWT.HandleFunc("/api/directories", servr.ServeDirectories)
	rJWT.HandleFunc("/api/import_status", servr.ServeImportStatus)
	rJWT.HandleFunc("/api/search", servr.ServeSearch)

	// begin api key routes
	rAPIKey := r.NewRoute().Subrouter()
	rAPIKey.Use(servr.WithJWTOrAPIKey())
	rAPIKey.HandleFunc("/api/upload", servr.ServeUpload)

	// frontend fallback route
	frontendFS := http.FS(frontend.FS)
	r.NotFoundHandler = http.FileServer(frontendFS)

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
