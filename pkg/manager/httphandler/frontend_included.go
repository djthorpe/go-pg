//go:build frontend

//go:generate sh -c "wasmbuild build -o frontend ../../../wasm/pgmanager && mv frontend/wasm_exec.html frontend/index.html"

package httphandler

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend/*
var frontendFS embed.FS

// RegisterFrontendHandler registers the frontend static file handler
func RegisterFrontendHandler(router *http.ServeMux, prefix string) {
	// Get the subdirectory to strip the "frontend" prefix
	subFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		panic(err)
	}

	// Serve static files from the embedded frontend folder
	router.Handle(joinPath(prefix, "/"), http.StripPrefix(prefix, http.FileServer(http.FS(subFS))))
}
