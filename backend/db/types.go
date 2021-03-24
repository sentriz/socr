package db

import "time"

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
	ID             int       `db:"id"                 json:"id"`
	Hash           string    `db:"hash"               json:"hash"`
	Timestamp      time.Time `db:"timestamp"          json:"timestamp"`
	DimWidth       int       `db:"dim_width"          json:"dim_width"`
	DimHeight      int       `db:"dim_height"         json:"dim_height"`
	DominantColour string    `db:"dominant_colour"    json:"dominant_colour"`
	Similarity     float64   `db:"similarity"         json:"similarity"`
	Blurhash       string    `db:"blurhash"           json:"blurhash"`
	Blocks         []*Block  `db:"highlighted_blocks" json:"highlighted_blocks"`
	Directories    []string  `db:"directories"        json:"directories"`
}
