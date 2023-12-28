package docs

import (
	"embed"
	_ "embed"

	rudidocs "go.xrstf.de/rudi/pkg/docs"
)

//go:embed *.md
var embeddedFS embed.FS

var Functions = rudidocs.NewFunctionProvider(&embeddedFS)
