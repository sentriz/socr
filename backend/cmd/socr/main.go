package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.senan.xyz/socr/controller"
	"go.senan.xyz/socr/imagery"

	_ "go.senan.xyz/socr/controller/auth"
)

const (
	screenshotIndex = "screenshots"
)

func mustEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("please provide a %q", key)
	return ""
}

func createIndexMapping() *mapping.IndexMappingImpl {
	fieldMapNumeric := bleve.NewNumericFieldMapping()

	fieldMapEnglish := bleve.NewTextFieldMapping()
	fieldMapEnglish.Analyzer = en.AnalyzerName

	// TODO: use this field mapping for tags
	fieldMapKeyword := bleve.NewTextFieldMapping()
	fieldMapKeyword.Analyzer = keyword.Name

	fieldMapTime := bleve.NewDateTimeFieldMapping()

	mappingBlocks := bleve.NewDocumentMapping()
	mappingBlocks.AddFieldMappingsAt("text", fieldMapEnglish)
	mappingBlocks.AddFieldMappingsAt("position", fieldMapNumeric)

	mappingScreenshot := bleve.NewDocumentMapping()
	mappingScreenshot.AddFieldMappingsAt("timestamp", fieldMapTime)
	mappingScreenshot.AddSubDocumentMapping("blocks", mappingBlocks)

	mappingIndex := bleve.NewIndexMapping()
	mappingIndex.DefaultMapping = mappingScreenshot

	return mappingIndex
}

func getOrCreateIndex(path string) (bleve.Index, error) {
	index, err := bleve.Open(path)
	switch {
	case
		errors.Is(err, bleve.ErrorIndexMetaMissing),
		errors.Is(err, bleve.ErrorIndexPathDoesNotExist):
		indexMapping := createIndexMapping()
		return bleve.New(path, indexMapping)
	case err != nil:
		return nil, fmt.Errorf("open index: %w", err)
	default:
		return index, nil
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	confListenAddr := mustEnv("SOCR_LISTEN_ADDR")
	confScreenshotsPath := mustEnv("SOCR_SCREENSHOTS_PATH")
	confIndexPath := mustEnv("SOCR_INDEX_PATH")
	confImportPath := mustEnv("SOCR_IMPORT_PATH")
	confHMACSecret := mustEnv("SOCR_HMAC_SECRET")
	confLoginUsername := mustEnv("SOCR_LOGIN_USERNAME")
	confLoginPassword := mustEnv("SOCR_LOGIN_PASSWORD")
	confAPIKey := mustEnv("SOCR_API_KEY")

	index, err := getOrCreateIndex(confIndexPath)
	if err != nil {
		log.Fatalf("error getting index: %v", err)
	}

	ctrl := &controller.Controller{
		ScreenshotsPath: confScreenshotsPath,
		ImportPath:      confImportPath,
		ImportRunning:   new(int32),
		Index:           index,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		SocketClientsSettings:   map[*websocket.Conn]struct{}{},
		SocketClientsScreenshot: map[string]map[*websocket.Conn]struct{}{},
		SocketUpdatesSettings:   make(chan struct{}),
		SocketUpdatesScreenshot: make(chan *controller.Screenshot),
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

	// begin authenticated routes
	rAuth := r.NewRoute().Subrouter()
	rAuth.Use(ctrl.WithAuth())
	rAuth.HandleFunc("/api/upload", ctrl.ServeUpload)
	rAuth.HandleFunc("/api/start_import", ctrl.ServeStartImport)
	rAuth.HandleFunc("/api/about", ctrl.ServeAbout)
	rAuth.HandleFunc("/api/import_status", ctrl.ServeImportStatus)
	rAuth.HandleFunc("/api/search", ctrl.ServeSearch)

	server := http.Server{
		Addr:    confListenAddr,
		Handler: r,
	}

	log.Printf("listening on %q", confListenAddr)
	log.Fatalf("error starting server: %v", server.ListenAndServe())
}
