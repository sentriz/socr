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
func NewScreenshot(screenshot db.SearchScreenshotsRow) (*Screenshot, error) {
	r := &Screenshot{
		ID:             hasher.ID(screenshot.ID),
		Timestamp:      screenshot.Timestamp,
		DirectoryAlias: screenshot.DirectoryAlias,
		Filename:       screenshot.Filename,
		DimWidth:       int(screenshot.DimWidth),
		DimHeight:      int(screenshot.DimHeight),
		DominantColour: screenshot.DominantColour,
		Blurhash:       screenshot.Blurhash,
		Blocks:         []*Block{},
	}

	var rawBlocks []*db.Block
	if err := json.Unmarshal(screenshot.Blocks, &rawBlocks); err != nil {
		return nil, fmt.Errorf("unmarshal blocks: %w", err)
	}
	for _, rawBlock := range rawBlocks {
		block, err := NewBlock(rawBlock)
		if err != nil {
			return nil, fmt.Errorf("convert block: %w", err)
		}
		r.Blocks = append(r.Blocks, block)
	}
	return r, nil
}

type Block struct {
	Index    int    `json:"index"`
	Position [4]int `json:"position"`
	Body     string `json:"body"`
}

func NewBlock(block *db.Block) (*Block, error) {
	return &Block{
		Index: int(block.Index),
		Position: [...]int{
			int(block.MinX),
			int(block.MinY),
			int(block.MaxX),
			int(block.MaxY),
		},
		Body: block.Body,
	}, nil
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
