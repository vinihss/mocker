package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
)

// SwaggerHandler handles Swagger UI endpoints
type SwaggerHandler struct {
	swaggerUIDir    string
	openAPISpecPath string
}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler(swaggerUIDir, openAPISpecPath string) *SwaggerHandler {
	return &SwaggerHandler{
		swaggerUIDir:    swaggerUIDir,
		openAPISpecPath: openAPISpecPath,
	}
}

// ServeSwaggerUI serves the Swagger UI index page
func (h *SwaggerHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(h.swaggerUIDir, "index.html")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		http.Error(w, "Swagger UI not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// ServeReDoc serves ReDoc documentation
func (h *SwaggerHandler) ServeReDoc(w http.ResponseWriter, r *http.Request) {
	reDocPath := filepath.Join(h.swaggerUIDir, "..", "redoc.html")

	data, err := os.ReadFile(reDocPath)
	if err != nil {
		http.Error(w, "ReDoc not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// ServeSwaggerJSON serves the OpenAPI spec as JSON
func (h *SwaggerHandler) ServeSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.openAPISpecPath)
	if err != nil {
		http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
		return
	}

	// Check if file is YAML (contains "openapi:" or "swagger:") or JSON
	content := strings.TrimSpace(string(data))

	if strings.HasPrefix(content, "openapi:") || strings.HasPrefix(content, "swagger:") {
		// It's YAML - convert to JSON
		var spec map[string]interface{}
		if err := yaml.Unmarshal(data, &spec); err != nil {
			http.Error(w, "Invalid OpenAPI spec", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(spec)
		if err != nil {
			http.Error(w, "Failed to convert spec to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	// Assume it's already JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// ServeStaticFile serves a static file from swagger-ui directory
func (h *SwaggerHandler) ServeStaticFile(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		http.Error(w, "File path required", http.StatusBadRequest)
		return
	}

	//Security: prevent directory traversal
	pathParam = strings.ReplaceAll(pathParam, "..", "")
	pathParam = strings.ReplaceAll(pathParam, "//", "/")

	fullPath := filepath.Join(h.swaggerUIDir, pathParam)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	contentType := getContentType(pathParam)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// ServeStaticFileFromChi serves static files using chi router (for /docs/* routes)
// Extracts the path from chi context and serves the file
func (h *SwaggerHandler) ServeStaticFileFromChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the wildcard path from chi
		path := chi.URLParam(r, "*")
		if path == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Security: prevent directory traversal
		path = strings.ReplaceAll(path, "..", "")
		path = strings.ReplaceAll(path, "//", "/")

		fullPath := filepath.Join(h.swaggerUIDir, path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		contentType := getContentType(path)
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})
}

// ServeStaticFileHandler returns an http.HandlerFunc for serving static files
func (h *SwaggerHandler) ServeStaticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path from chi URL
		path := chi.URLParam(r, "*")
		if path == "" {
			http.Error(w, "File path required", http.StatusBadRequest)
			return
		}

		// Security: prevent directory traversal
		path = strings.ReplaceAll(path, "..", "")
		path = strings.ReplaceAll(path, "//", "/")

		fullPath := filepath.Join(h.swaggerUIDir, path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		contentType := getContentType(path)
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

func getContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

// ServeFile serves a single file from the given directory
func ServeFile(w http.ResponseWriter, r *http.Request, dir, filename string) {
	fullPath := filepath.Join(dir, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	contentType := getContentType(filename)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// ReadFile reads a file from the given path
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
