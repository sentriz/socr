package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"go.senan.xyz/flagconf"
	"golang.org/x/sync/errgroup"

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
	confScanInterval := flag.Duration("scan-interval", 0, "interval between periodic scans of source directories")

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbc, err := db.New(ctx, *confDBDSN)
	if err != nil {
		log.Panicf("error creating database: %v", err)
	}
	defer dbc.Close()

	if err := dbc.Migrate(ctx); err != nil {
		log.Panicf("error running migrations: %v", err)
	}

	importr := importer.New(dbc, png.Encode, "image/png", directories.Directories(confDirs), *confUploadsAlias, *confThumbnailWidth)
	servr := server.New(dbc, importr, directories.Directories(confDirs), *confUploadsAlias, *confHMACSecret, *confLoginUsername, *confLoginPassword, *confAPIKey)

	errgrp, ctx := errgroup.WithContext(ctx)

	numImportWorkers := runtime.NumCPU()
	for i := range numImportWorkers {
		errgrp.Go(func() error {
			defer logJob("import worker", "n", i+1)()
			return importr.StartWorker(ctx)
		})
	}

	errgrp.Go(func() error {
		defer logJob("watch updates")()
		return importr.WatchUpdates(ctx)
	})

	errgrp.Go(func() error {
		defer logJob("scan loop")()
		return importr.RunScanLoop(ctx)
	})

	if *confScanInterval > 0 {
		errgrp.Go(func() error {
			defer logJob("periodic scan", "interval", *confScanInterval)()
			return importr.RunPeriodicScan(ctx, *confScanInterval)
		})
	}

	errgrp.Go(func() error {
		defer logJob("socket notify scanner update")()
		return servr.SocketNotifyScannerUpdate(ctx)
	})

	errgrp.Go(func() error {
		defer logJob("socket notify media")()
		return servr.SocketNotifyMedia(ctx)
	})

	errgrp.Go(func() error {
		defer logJob("http", "addr", *confListenAddr)()

		httpServer := &http.Server{
			Addr:              *confListenAddr,
			Handler:           servr.Router(),
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1024 * 64,
			BaseContext:       func(l net.Listener) context.Context { return ctx },
		}
		errgrp.Go(func() error {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
			defer cancel()
			return httpServer.Shutdown(shutdownCtx)
		})

		log.Printf("starting socr %s", socr.Version)
		log.Printf("listening on %q", *confListenAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	if err := errgrp.Wait(); err != nil {
		log.Panic(err)
	}

	log.Print("shutdown complete")
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
		return errors.New("alias and path must be non-empty")
	}
	d[alias] = path
	return nil
}

func logJob(jobName string, args ...any) func() {
	log.Printf("starting job %q %v", jobName, args)
	return func() { log.Printf("stopped job %q", jobName) }
}
