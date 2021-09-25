//nolint:gochecknoglobals,golint,stylecheck
package socr

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed version.txt
var version string
var Version = fmt.Sprintf("v%s", strings.TrimSpace(version))
