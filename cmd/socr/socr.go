//nolint:gochecknoglobals
package main

import (
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.senan.xyz/socr"
	"go.senan.xyz/socr/pkg/db"
	"go.senan.xyz/socr/pkg/directories"
	"go.senan.xyz/socr/pkg/importer"
	"go.senan.xyz/socr/pkg/server"
)

var (
	confListenAddr     = mustEnv("SOCR_LISTEN_ADDR")
	confDBDSN          = mustEnv("SOCR_DB_DSN")
	confHMACSecret     = mustEnv("SOCR_HMAC_SECRET")
	confLoginUsername  = mustEnv("SOCR_LOGIN_USERNAME")
	confLoginPassword  = mustEnv("SOCR_LOGIN_PASSWORD")
	confAPIKey         = mustEnv("SOCR_API_KEY")
	confDirs           = envDirs("SOCR_DIR_")
	confUploadsAlias   = envOr("SOCR_UPLOADS_DIR_ALIAS", "uploads")
	confThumbnailWidth = envOrInt("SOCR_THUMBNAIL_WIDTH", 315)
)

func main() {
	if _, ok := confDirs[confUploadsAlias]; !ok {
		log.Fatalf("please provide an uploads directory")
	}
	for alias, path := range confDirs {
		log.Printf("using directory alias %q path %q", alias, path)
	}

	dbc, err := db.New(confDBDSN)
	if err != nil {
		log.Panicf("error creating database: %v", err)
	}
	defer dbc.Close()

	if err := dbc.Migrate(); err != nil {
		log.Panicf("error running migrations: %v", err)
	}

	const numImportWorkers = 1
	importr := importer.New(dbc, png.Encode, "image/png", confDirs, confUploadsAlias, uint(confThumbnailWidth))
	for i := 0; i < numImportWorkers; i++ {
		log.Printf("starting import worker %d", i+1)
		go importr.StartWorker()
	}
	go func() {
		if err := importr.WatchUpdates(); err != nil {
			log.Printf("error starting watcher: %v", err)
		}
	}()

	servr := server.New(dbc, importr, confDirs, confUploadsAlias, confHMACSecret, confLoginUsername, confLoginPassword, confAPIKey)
	go servr.SocketNotifyScannerUpdate()
	go servr.SocketNotifyMedia()

	router := servr.Router()
	server := http.Server{
		Addr:              confListenAddr,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1024 * 64,
	}

	log.Printf("starting socr %s", socr.Version)
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

func envOr(key string, or string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return or
}

func envOrInt(key string, or int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return or
}

func envDirs(prefix string) directories.Directories {
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
