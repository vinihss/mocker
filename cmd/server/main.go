package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vinicius/mocker/internal/config"
	"github.com/vinicius/mocker/internal/handlers"
	"github.com/vinicius/mocker/internal/middleware"
	"github.com/vinicius/mocker/internal/storage"
)

// Get executable directory for finding swagger-ui
func getExecDir() string {
	execPath, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(execPath)
}

// getSwaggerUIDir returns the directory containing Swagger UI files
func getSwaggerUIDir() string {
	// Try multiple locations for swagger-ui
	locations := []string{
		"./docs/swagger-ui",
		"docs/swagger-ui",
		filepath.Join(getExecDir(), "docs/swagger-ui"),
		filepath.Join(getExecDir(), "./docs/swagger-ui"),
		"/home/vinicius/mocker/docs/swagger-ui",
	}

	for _, loc := range locations {
		if _, err := os.Stat(filepath.Join(loc, "index.html")); err == nil {
			return loc
		}
	}

	// Default to local path
	return "./docs/swagger-ui"
}

// getOpenAPISpecPath returns the path to the OpenAPI spec file
func getOpenAPISpecPath() string {
	locations := []string{
		"./docs/openapi.yaml",
		"docs/openapi.yaml",
		filepath.Join(getExecDir(), "docs/openapi.yaml"),
		filepath.Join(getExecDir(), "./docs/openapi.yaml"),
		"/home/vinicius/mocker/docs/openapi.yaml",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Default to local path
	return "./docs/openapi.yaml"
}

func main() {
	// Load configuration
	cfg := config.Default()

	// Initialize database
	store, err := storage.NewSQLite(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	// Initialize handlers
	h := handlers.New(store)

	// Setup router
	r := setupRouter(h)

	// Start server
	addr := cfg.Server.GetPort()
	log.Printf("Starting server on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Println("Server started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.CORS)

	// Swagger UI - determine paths
	swaggerUIDir := getSwaggerUIDir()
	openAPISpecPath := getOpenAPISpecPath()
	swaggerHandler := handlers.NewSwaggerHandler(swaggerUIDir, openAPISpecPath)

	// Swagger UI routes
	r.Route("/docs", func(r chi.Router) {
		r.Get("/", swaggerHandler.ServeSwaggerUI)
		r.Get("/redoc", swaggerHandler.ServeReDoc)
		r.Get("/swagger.json", swaggerHandler.ServeSwaggerJSON)
		r.Get("/*", swaggerHandler.ServeStaticFileHandler())
	})

	// Health check and frontend routes separated
	r.Get("/health", h.HealthCheck)

	// Mock server - serve mocks directly on their configured paths
	r.MethodFunc("GET", "/api/users", h.ServeMock)
	r.MethodFunc("POST", "/api/users", h.ServeMock)
	r.MethodFunc("PUT", "/api/users", h.ServeMock)
	r.MethodFunc("DELETE", "/api/users", h.ServeMock)
	r.MethodFunc("PATCH", "/api/users", h.ServeMock)
	r.MethodFunc("GET", "/api/products", h.ServeMock)
	r.MethodFunc("POST", "/api/products", h.ServeMock)
	r.MethodFunc("PUT", "/api/products", h.ServeMock)
	r.MethodFunc("DELETE", "/api/products", h.ServeMock)
	r.MethodFunc("GET", "/api/users/{id}", h.ServeMock)
	r.MethodFunc("PUT", "/api/users/{id}", h.ServeMock)
	r.MethodFunc("DELETE", "/api/users/{id}", h.ServeMock)
	r.MethodFunc("PATCH", "/api/users/{id}", h.ServeMock)
	r.MethodFunc("GET", "/api/products/{id}", h.ServeMock)
	r.MethodFunc("PUT", "/api/products/{id}", h.ServeMock)
	r.MethodFunc("DELETE", "/api/products/{id}", h.ServeMock)
	r.MethodFunc("PATCH", "/api/products/{id}", h.ServeMock)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Mock routes
		r.Route("/mocks", func(r chi.Router) {
			r.Post("/", h.CreateMock)
			r.Get("/", h.ListMocks)
			r.Get("/{id}", h.GetMock)
			r.Put("/{id}", h.UpdateMock)
			r.Delete("/{id}", h.DeleteMock)
			r.Post("/{id}/test", h.TestMock)
			r.Get("/{id}/tests", h.GetMockTests)
			r.Post("/{id}/activate", h.ActivateMock)
			r.Post("/{id}/deactivate", h.DeactivateMock)
		})
	})

	// 404 handler
	r.NotFound(h.NotFound)

	return r
}
