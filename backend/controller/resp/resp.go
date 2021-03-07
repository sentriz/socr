package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
)

func Write(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Response interface{} `json:"result"`
	}{
		Response: body,
	})
}

func Error(w http.ResponseWriter, status int, format string, a ...interface{}) {
	w.WriteHeader(status)
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, a...),
	})
}

type ID struct {
	ID hasher.ID `json:"id"`
}

type Status struct {
	Status string `json:"status"`
}

type Block struct {
	Index    int    `json:"index"`
	Position [4]int `json:"position"`
	Body     string `json:"body"`
}

func NewBlock(dbBlock *db.Blocks) *Block {
	return &Block{
		Index: int(dbBlock.Index.Int),
		Position: [...]int{
			int(dbBlock.MinX.Int),
			int(dbBlock.MinY.Int),
			int(dbBlock.MaxX.Int),
			int(dbBlock.MaxY.Int),
		},
		Body: dbBlock.Body.String,
	}
}

type Screenshot struct {
	ID             hasher.ID `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	DirectoryAlias string    `json:"directory_alias"`
	Filename       string    `json:"filename"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
	Blocks         []*Block  `json:"blocks"`
}

// TODO: try use db.Screenshot here instead of db.SearchScreenshotsRow
// perhaps https://github.com/kyleconroy/sqlc/issues/755 ?
// perhaps https://github.com/kyleconroy/sqlc/issues/643 ?
func NewScreenshot(dbScreenshot db.SearchScreenshotsRow) *Screenshot {
	screenshot := &Screenshot{
		ID:             hasher.ID(dbScreenshot.ID.Int),
		Timestamp:      dbScreenshot.Timestamp.Time,
		DirectoryAlias: dbScreenshot.DirectoryAlias.String,
		Filename:       dbScreenshot.Filename.String,
		DimWidth:       int(dbScreenshot.DimWidth.Int),
		DimHeight:      int(dbScreenshot.DimHeight.Int),
		DominantColour: dbScreenshot.DominantColour.String,
		Blurhash:       dbScreenshot.Blurhash.String,
	}
	for _, dbBlock := range dbScreenshot.Blocks {
		screenshot.Blocks = append(screenshot.Blocks, NewBlock(&dbBlock))
	}
	return screenshot
}

func NewScreenshots(dbScreenshots []db.SearchScreenshotsRow) []*Screenshot {
	var screenshots []*Screenshot
	for _, dbScreenshot := range dbScreenshots {
		screenshots = append(screenshots, NewScreenshot(dbScreenshot))
	}
	return screenshots
}

type ScreenshotCount struct {
	DirectoryAlias string
	Count          int
}

func NewScreenshotCount(dbScreenshotCount db.CountDirectoriesByAliasRow) *ScreenshotCount {
	return &ScreenshotCount{
		DirectoryAlias: dbScreenshotCount.DirectoryAlias.String,
		Count:          int(dbScreenshotCount.Count.Int),
	}
}

type About struct {
	Version          string                          `json:"version"`
	APIKey           string                          `json:"api_key"`
	SocketClients    int                             `json:"socket_clients"`
	ImportPath       string                          `json:"import_path"`
	ScreenshotsPath  string                          `json:"screenshots_path"`
	ScreenshotsCount []db.CountDirectoriesByAliasRow `json:"screenshots_indexed"`
}

type Token struct {
	Token string `json:"token"`
}

type ImportStatus struct {
	Running bool `json:"running"`
}
