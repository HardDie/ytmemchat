package chat

import (
	_ "embed"
	"net/http"
)

var (
	//go:embed index.html
	overlayHTML []byte
)

// HTMLHandler serves the main overlay.html file when OBS requests the root path.
func (*Chat) HTMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat" {
		// This should be handled by the file server, but for robustness:
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Serve overlay.html from the current directory
	w.Write(overlayHTML)
}
