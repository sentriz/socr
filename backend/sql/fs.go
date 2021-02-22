package sql

import (
	_ "embed"
)

//go:embed schema.pgsql
var Schema string
