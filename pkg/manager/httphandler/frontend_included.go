//go:build frontend

//go:generate sh -c "wasmbuild build -o frontend ../../../wasm/pgmanager && mv frontend/wasm_exec.html frontend/index.html"

package httphandler

import (
	"embed"
	"net/http"
)

//go:embed frontend/*
var frontendFS embed.FS

// RegisterFrontendHandler registers the frontend static file handler
func RegisterFrontendHandler(router *http.ServeMux, prefix string) {
	// Serve static files from the embedded frontend folder
	fileServer := http.FileServer(http.FS(frontendFS))
	router.Handle(joinPath(prefix, "/"), http.StripPrefix(prefix, fileServer))
}
