package docs

import (
	"embed"
	_ "embed"
)

//go:embed *.md
var embeddedFS embed.FS

var Functions = helpProvider{}

type helpProvider struct{}

func (fd helpProvider) Documentation(functionName string) (string, error) {
	contents, err := embeddedFS.ReadFile(functionName + ".md")
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
