package assets

import (
	_ "embed"
)

//go:embed templates/style.css
var PageCSS string

//go:embed templates/page.html
var PageHTML string
