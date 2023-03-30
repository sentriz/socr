//nolint:gochecknoglobals
package web

import (
	"embed"
	"io/fs"
)

//go:generate npm install
//go:generate npm run-script build

//go:embed dist
var dist embed.FS
var Dist, _ = fs.Sub(dist, "dist")

//go:embed dist/index.html
var Index []byte
