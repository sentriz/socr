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
func NewScreenshot(dbScreenshot db.SearchScreenshotsRow) (*Screenshot, error) {
	screenshot := &Screenshot{
		ID:             hasher.ID(dbScreenshot.ID),
		Timestamp:      dbScreenshot.Timestamp,
		DirectoryAlias: dbScreenshot.DirectoryAlias,
		Filename:       dbScreenshot.Filename,
		DimWidth:       int(dbScreenshot.DimWidth),
		DimHeight:      int(dbScreenshot.DimHeight),
		DominantColour: dbScreenshot.DominantColour,
		Blurhash:       dbScreenshot.Blurhash,
	}

	var dbBlocks []*db.Block
	if err := json.Unmarshal(dbScreenshot.Blocks, &dbBlocks); err != nil {
		return nil, fmt.Errorf("unmarshal blocks: %w", err)
	}
	for _, dbBlock := range dbBlocks {
		block, err := NewBlock(dbBlock)
		if err != nil {
			return nil, fmt.Errorf("convert block: %w", err)
		}
		screenshot.Blocks = append(screenshot.Blocks, block)
	}
	return screenshot, nil
}

type Block struct {
	Index    int    `json:"index"`
	Position [4]int `json:"position"`
	Body     string `json:"body"`
}

func NewBlock(dbBlock *db.Block) (*Block, error) {
	return &Block{
		Index: int(dbBlock.Index),
		Position: [...]int{
			int(dbBlock.MinX),
			int(dbBlock.MinY),
			int(dbBlock.MaxX),
			int(dbBlock.MaxY),
		},
		Body: dbBlock.Body,
	}, nil
}

func NewScreenshots(dbScreenshots []db.SearchScreenshotsRow) ([]*Screenshot, error) {
	var screenshots []*Screenshot
	for _, dbScreenshot := range dbScreenshots {
		screenshot, err := NewScreenshot(dbScreenshot)
		if err != nil {
			return nil, fmt.Errorf("convert screenshot: %w", err)
		}
		screenshots = append(screenshots, screenshot)
	}
	return screenshots, nil
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
