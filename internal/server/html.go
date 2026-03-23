package server

import (
	_ "embed"
	"net/http"
)

var (
	//go:embed overlay.html
	overlayHTML []byte
	//go:embed favicon.png
	favicon []byte
)

// htmlHandler serves the main overlay.html file when OBS requests the root path.
func htmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// This should be handled by the file server, but for robustness:
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Serve overlay.html from the current directory
	w.Write(overlayHTML)
}

// faviconHandler serves the favicon.png file when OBS requests the /favicon.ico path.
func faviconHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	// Serve overlay.html from the current directory
	w.Write(favicon)
}
