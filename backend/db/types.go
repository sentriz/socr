package db

import "time"

type Block struct {
	ID           int    `db:"id"`
	ScreenshotID int    `db:"screenshot_id"`
	Index        int    `db:"index"`
	MinX         int    `db:"min_x"`
	MinY         int    `db:"min_y"`
	MaxX         int    `db:"max_x"`
	MaxY         int    `db:"max_y"`
	Body         string `db:"body"`
}

type Screenshot struct {
	ID             int       `db:"id"`
	Hash           string    `db:"hash"`
	Timestamp      time.Time `db:"timestamp"`
	DimWidth       int       `db:"dim_width"`
	DimHeight      int       `db:"dim_height"`
	DominantColour string    `db:"dominant_colour"`
	Similarity     float64   `db:"similarity"`
	Blurhash       string    `db:"blurhash"`
	Blocks         []*Block  `db:"highlighted_blocks"`
	Directories    []string  `db:"directories"`
}
