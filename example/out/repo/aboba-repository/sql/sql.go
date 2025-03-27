package sql

import (
	_ "embed"
)

//go:embed create.sql
var Create string

//go:embed get.sql
var Get string
