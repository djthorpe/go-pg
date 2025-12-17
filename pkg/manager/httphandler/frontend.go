//go:build frontend

package httphandler

import (
	"net/http"

	// Packages
	frontend "github.com/mutablelogic/go-pg/build/wasm/pgmanager"
)

// RegisterFrontendHandler registers the frontend static file handler
func RegisterFrontendHandler(router *http.ServeMux, prefix string) {
	// Serve static files
	fileServer := http.FileServer(http.FS(frontend.FS))
	router.Handle(joinPath(prefix, "/"), http.StripPrefix(prefix, fileServer))
}
