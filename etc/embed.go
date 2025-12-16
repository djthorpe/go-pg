package frontend

import "embed"

//go:embed *.html *.js *.wasm *.png
var FS embed.FS
