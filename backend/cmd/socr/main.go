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
	bleveHTTP "github.com/blevesearch/bleve/http"
	"github.com/blevesearch/bleve/mapping"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.senan.xyz/socr/controller"

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
	confPassword := mustEnv("SOCR_PASSWORD")
	confAPIKey := mustEnv("SOCR_API_KEY")

	index, err := getOrCreateIndex(confIndexPath)
	if err != nil {
		log.Fatalf("error getting index: %v", err)
	}

	ctrl := &controller.Controller{
		ScreenshotsPath: confScreenshotsPath,
		ImportPath:      confImportPath,
		ImportUpdates:   make(chan controller.ImportUpdate),
		Index:           index,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		SocketClients: map[*websocket.Conn]struct{}{},
		HMACSecret:    confHMACSecret,
		Password:      confPassword,
		APIKey:        confAPIKey,
	}

	go ctrl.EmitImportUpdates()

	r := mux.NewRouter()
	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"DNT", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Range"}),
		handlers.MaxAge(1728000),
	))

	r.HandleFunc("/api/upload", ctrl.ServeUpload)
	r.HandleFunc("/api/image/{id}", ctrl.ServeImage)
	r.HandleFunc("/api/start_import", ctrl.ServeStartImport)
	r.HandleFunc("/api/ws", ctrl.ServeWebSocket)
	r.HandleFunc("/api/authenticate", ctrl.ServeAuthenticate)

	bleveHTTP.RegisterIndexName(screenshotIndex, index)
	r.Handle("/api/search", bleveHTTP.NewSearchHandler(screenshotIndex))

	server := http.Server{
		Addr:    confListenAddr,
		Handler: r,
	}

	log.Printf("listening on %q", confListenAddr)
	log.Fatalf("error starting server: %v", server.ListenAndServe())
}
