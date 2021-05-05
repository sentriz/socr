package socr

import (
	"embed"
	"io/fs"
)

//go:embed dist/index.html
var Index []byte

//go:embed dist/favicon.ico
var Favicon []byte

//go:embed dist/assets
var assets embed.FS
var Assets, _ = fs.Sub(assets, "dist")
