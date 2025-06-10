package templates

import "embed"

//go:embed index.html index.css index.js
var TemplateFS embed.FS
