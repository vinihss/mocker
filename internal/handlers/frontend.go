package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FrontendHandler serves the React frontend static files
type FrontendHandler struct {
	frontendDir string
}

// NewFrontendHandler creates a new Frontend handler
func NewFrontendHandler(frontendDir string) *FrontendHandler {
	return &FrontendHandler{
		frontendDir: frontendDir,
	}
}

// ServeFrontend serves the React frontend
func (h *FrontendHandler) ServeFrontend(w http.ResponseWriter, r *http.Request) {
	// For SPA, serve index.html for all routes that don't have a file extension
	path := r.URL.Path

	// Check if path has a file extension
	if !strings.Contains(filepath.Ext(path), ".") && !strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".css") {
		// SPA - serve index.html
		indexPath := filepath.Join(h.frontendDir, "index.html")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			http.Error(w, "Frontend not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
		return
	}

	// Serve static file
	fullPath := filepath.Join(h.frontendDir, path)

	// Security: prevent directory traversal
	fullPath = strings.ReplaceAll(fullPath, "..", "")

	data, err := os.ReadFile(fullPath)
	if err != nil {
		// Fallback to index.html for SPA
		indexPath := filepath.Join(h.frontendDir, "index.html")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			http.Error(w, "Frontend not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
		return
	}

	contentType := getContentType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// ServeStatic serves static files from frontend directory
func (h *FrontendHandler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Remove leading slash
	path = strings.TrimPrefix(path, "/")

	fullPath := filepath.Join(h.frontendDir, path)

	// Security: prevent directory traversal
	fullPath = strings.ReplaceAll(fullPath, "..", "")

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	contentType := getContentType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
