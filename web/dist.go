//nolint:gochecknoglobals
package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS
var Dist, _ = fs.Sub(dist, "dist")
