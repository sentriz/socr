//nolint:gochecknoglobals
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/directories"
	"go.senan.xyz/socr/backend/imagery"
	"go.senan.xyz/socr/backend/importer"
	"go.senan.xyz/socr/backend/scanner"
	"go.senan.xyz/socr/backend/server"
)

const uploadsAlias = "uploads"

var (
	confListenAddr    = mustEnv("SOCR_LISTEN_ADDR")
	confDBDSN         = mustEnv("SOCR_DB_DSN")
	confHMACSecret    = mustEnv("SOCR_HMAC_SECRET")
	confLoginUsername = mustEnv("SOCR_LOGIN_USERNAME")
	confLoginPassword = mustEnv("SOCR_LOGIN_PASSWORD")
	confAPIKey        = mustEnv("SOCR_API_KEY")
	confDirs          = parseEnvDirs("SOCR_DIR_")
)

func main() {
	if _, ok := confDirs[uploadsAlias]; !ok {
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
		Running:                 new(int32),
		Directories:             confDirs,
		DirectoriesUploadsAlias: uploadsAlias,
		DB:                      dbConn,
		Importer:                importr,
		Updates:                 make(chan struct{}),
	}
	servr := &server.Server{
		Directories:             confDirs,
		DirectoriesUploadsAlias: uploadsAlias,
		DB:                      dbConn,
		Importer:                importr,
		Scanner:                 scanr,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
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
	go func() {
		if err := scanr.WatchUpdates(); err != nil {
			log.Printf("error starting watcher: %v", err)
		}
	}()

	router := servr.Router()
	server := http.Server{
		Addr:    confListenAddr,
		Handler: router,
	}

	log.Printf("listening on %q", confListenAddr)
	log.Printf("starting server: %v", server.ListenAndServe())
}

func mustEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("please provide a %q", key)
	return ""
}

func parseEnvDirs(prefix string) directories.Directories {
	expr := regexp.MustCompile(prefix + `(?P<Alias>[\w_]+)=(?P<Path>.*)`)
	const (
		partFull = iota
		partAlias
		partPath
	)
	dirMap := directories.Directories{}
	for _, env := range os.Environ() {
		parts := expr.FindStringSubmatch(env)
		if len(parts) != 3 {
			continue
		}
		alias := strings.ToLower(parts[partAlias])
		path := filepath.Clean(parts[partPath])
		dirMap[alias] = path
	}
	return dirMap
}
