package db

import (
	"database/sql"
	"time"
)

type Block struct {
	ID           int    `db:"id"            json:"id"`
	ScreenshotID int    `db:"screenshot_id" json:"screenshot_id"`
	Index        int    `db:"index"         json:"index"`
	MinX         int    `db:"min_x"         json:"min_x"`
	MinY         int    `db:"min_y"         json:"min_y"`
	MaxX         int    `db:"max_x"         json:"max_x"`
	MaxY         int    `db:"max_y"         json:"max_y"`
	Body         string `db:"body"          json:"body"`
}

type Screenshot struct {
	ID                int       `db:"id"                 json:"id"`
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
	ID             int           `db:"id"              json:"id"`
	ScreenshotID   sql.NullInt64 `sql:"screenshot_id"   json:"screenshot_id,omitempty"`
	VideoID        sql.NullInt64 `sql:"video_id"        json:"video_id,omitempty"`
	Filename       string        `db:"filename"        json:"filename"`
	DirectoryAlias string        `db:"directory_alias" json:"directory_alias"`
}

type DirectoryCount struct {
	DirectoryAlias string `db:"directory_alias" json:"directory_alias"`
	Count          int    `db:"count"           json:"count"`
}
