package db

import (
	"time"
)

type Block struct {
	ID      int    `db:"id"       json:"id"`
	MediaID int    `db:"media_id" json:"media_id"`
	Index   int    `db:"index"    json:"index"`
	MinX    int    `db:"min_x"    json:"min_x"`
	MinY    int    `db:"min_y"    json:"min_y"`
	MaxX    int    `db:"max_x"    json:"max_x"`
	MaxY    int    `db:"max_y"    json:"max_y"`
	Body    string `db:"body"     json:"body"`
}

type MediaType string

const (
	MediaTypeScreenshot MediaType = "screenshot"
	MediaTypeVideo      MediaType = "video"
)

type Media struct {
	ID                int       `db:"id"                 json:"id"`
	Type              MediaType `db:"type"               json:"type"`
	Hash              string    `db:"hash"               json:"hash"`
	Timestamp         time.Time `db:"timestamp"          json:"timestamp"`
	DimWidth          int       `db:"dim_width"          json:"dim_width"`
	DimHeight         int       `db:"dim_height"         json:"dim_height"`
	DominantColour    string    `db:"dominant_colour"    json:"dominant_colour"`
	Blurhash          string    `db:"blurhash"           json:"blurhash"`
	Similarity        float64   `db:"similarity"         json:"similarity,omitempty"`
	Blocks            []*Block  `db:"blocks"             json:"blocks,omitempty"`
	HighlightedBlocks []*Block  `db:"highlighted_blocks" json:"highlighted_blocks,omitempty"`
	Directories       []string  `db:"directories"        json:"directories,omitempty"`
}

type DirInfo struct {
	MediaID        int    `db:"media_id"        json:"media_id"`
	Filename       string `db:"filename"        json:"filename"`
	DirectoryAlias string `db:"directory_alias" json:"directory_alias"`
}

type DirectoryCount struct {
	DirectoryAlias string `db:"directory_alias" json:"directory_alias"`
	Count          int    `db:"count"           json:"count"`
}
