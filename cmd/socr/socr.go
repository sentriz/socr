package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go.senan.xyz/flagconf"
	"go.senan.xyz/socr"
	"go.senan.xyz/socr/db"
	"go.senan.xyz/socr/directories"
	"go.senan.xyz/socr/importer"
	"go.senan.xyz/socr/server"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

func main() {
	confListenAddr := flag.String("listen-addr", "", "address to listen on")
	confDBDSN := flag.String("db-dsn", "", "postgres connection string")
	confHMACSecret := flag.String("hmac-secret", "", "secret used to sign tokens")
	confLoginUsername := flag.String("login-username", "", "login username")
	confLoginPassword := flag.String("login-password", "", "login password")
	confAPIKey := flag.String("api-key", "", "api key")
	confUploadsAlias := flag.String("uploads-dir-alias", "uploads", "alias of the uploads directory")
	confThumbnailWidth := flag.Uint("thumbnail-width", 315, "thumbnail width in pixels")

	var confDirs = dirsFlag{}
	flag.Var(&confDirs, "dir", "directory in the form alias=path (repeatable)")

	confConfigPath := flag.String("config-path", "", "path to config file")

	flag.Parse()
	flagconf.ParseEnv()
	flagconf.ParseConfig(*confConfigPath)

	if *confListenAddr == "" {
		log.Fatalf("please provide a listen-addr")
	}
	if *confDBDSN == "" {
		log.Fatalf("please provide a db-dsn")
	}
	if *confHMACSecret == "" {
		log.Fatalf("please provide a hmac-secret")
	}
	if *confLoginUsername == "" {
		log.Fatalf("please provide a login-username")
	}
	if *confLoginPassword == "" {
		log.Fatalf("please provide a login-password")
	}
	if *confAPIKey == "" {
		log.Fatalf("please provide a api-key")
	}

	if _, ok := confDirs[*confUploadsAlias]; !ok {
		log.Fatalf("please provide an uploads directory")
	}
	for alias, path := range confDirs {
		log.Printf("using directory alias %q path %q", alias, path)
	}

	dbc, err := db.New(*confDBDSN)
	if err != nil {
		log.Panicf("error creating database: %v", err)
	}
	defer dbc.Close()

	if err := dbc.Migrate(); err != nil {
		log.Panicf("error running migrations: %v", err)
	}

	const numImportWorkers = 1
	importr := importer.New(dbc, png.Encode, "image/png", directories.Directories(confDirs), *confUploadsAlias, *confThumbnailWidth)
	for i := range numImportWorkers {
		log.Printf("starting import worker %d", i+1)
		go importr.StartWorker()
	}
	go func() {
		if err := importr.WatchUpdates(); err != nil {
			log.Printf("error starting watcher: %v", err)
		}
	}()

	servr := server.New(dbc, importr, directories.Directories(confDirs), *confUploadsAlias, *confHMACSecret, *confLoginUsername, *confLoginPassword, *confAPIKey)
	go servr.SocketNotifyScannerUpdate()
	go servr.SocketNotifyMedia()

	router := servr.Router()
	server := http.Server{
		Addr:              *confListenAddr,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1024 * 64,
	}

	log.Printf("starting socr %s", socr.Version)
	log.Printf("listening on %q", *confListenAddr)
	log.Printf("starting server %v", server.ListenAndServe())
}

type dirsFlag directories.Directories

func (d dirsFlag) String() string {
	var parts []string
	for alias, path := range d {
		parts = append(parts, alias+"="+path)
	}
	return strings.Join(parts, ", ")
}

func (d dirsFlag) Set(v string) error {
	alias, path, ok := strings.Cut(v, "=")
	if !ok {
		return fmt.Errorf("expected alias=path, got %q", v)
	}
	alias = strings.ToLower(strings.TrimSpace(alias))
	path = filepath.Clean(strings.TrimSpace(path))
	if alias == "" || path == "" {
		return fmt.Errorf("alias and path must be non-empty")
	}
	d[alias] = path
	return nil
}
