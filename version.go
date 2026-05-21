//nolint:gochecknoglobals
package socr

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var version string
var Version = "v" + strings.TrimSpace(version)
