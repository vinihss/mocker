package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/vinicius/mocker/internal/generator"
	"github.com/vinicius/mocker/internal/models"
	"github.com/vinicius/mocker/internal/storage"
	"github.com/vinicius/mocker/internal/validator"
)

// Handler holds all handlers
type Handler struct {
	store *storage.SQLite
}

// New creates a new handler
func New(store *storage.SQLite) *Handler {
	return &Handler{
		store: store,
	}
}

// Mock handlers

// CreateMock handles POST /api/v1/mocks
func (h *Handler) CreateMock(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if err := validator.Validate(req.Name, req.Path, req.Method); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate response config
	if req.ResponseConfig.Type == "faker" && req.ResponseConfig.Faker == nil {
		h.sendError(w, http.StatusBadRequest, "responseConfig.faker is required for faker type")
		return
	}
	if req.ResponseConfig.Type != "faker" && req.ResponseConfig.Body == nil {
		h.sendError(w, http.StatusBadRequest, "responseConfig.body is required for static/dynamic type")
		return
	}
	if req.ResponseConfig.StatusCode == 0 {
		req.ResponseConfig.StatusCode = 200
	}
	if err := validator.ValidateStatusCode(req.ResponseConfig.StatusCode); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	mock := &models.Mock{
		ID:              uuid.New().String(),
		Name:            req.Name,
		Path:            req.Path,
		Method:          req.Method,
		Description:     req.Description,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
		ResponseType:    req.ResponseConfig.Type,
		ResponseStatus:  req.ResponseConfig.StatusCode,
		ResponseDelayMs: req.ResponseConfig.DelayMs,
	}

	// Convert ResponseConfig to DB fields
	if err := mock.FromResponseConfig(req.ResponseConfig); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid response config")
		return
	}

	if err := h.store.CreateMock(r.Context(), mock); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert back to response
	_ = mock.ToResponseConfig()
	h.sendJSON(w, http.StatusCreated, mock)
}

// ListMocks handles GET /api/v1/mocks
func (h *Handler) ListMocks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	mocks, total, err := h.store.ListMocks(r.Context(), page, pageSize)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert to response format
	response := make([]models.Mock, len(mocks))
	for i, m := range mocks {
		_ = m.ToResponseConfig()
		response[i] = m
	}

	totalPages := (total + pageSize - 1) / pageSize
	h.sendJSON(w, http.StatusOK, models.ListResponse{
		Data:       response,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetMock handles GET /api/v1/mocks/:id
func (h *Handler) GetMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mock, err := h.store.GetMockByID(r.Context(), id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Mock not found")
		return
	}

	_ = mock.ToResponseConfig()
	h.sendJSON(w, http.StatusOK, mock)
}

// UpdateMock handles PUT /api/v1/mocks/:id
func (h *Handler) UpdateMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get existing mock
	mock, err := h.store.GetMockByID(r.Context(), id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Mock not found")
		return
	}

	var req models.UpdateMockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Apply updates
	if req.Name != nil {
		mock.Name = *req.Name
	}
	if req.Path != nil {
		mock.Path = *req.Path
	}
	if req.Method != nil {
		mock.Method = *req.Method
	}
	if req.Description != nil {
		mock.Description = *req.Description
	}
	if req.IsActive != nil {
		mock.IsActive = *req.IsActive
	}
	if req.ResponseConfig != nil {
		if err := mock.FromResponseConfig(*req.ResponseConfig); err != nil {
			h.sendError(w, http.StatusBadRequest, "Invalid response config")
			return
		}
	}

	mock.UpdatedAt = time.Now()

	if err := h.store.UpdateMock(r.Context(), mock); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = mock.ToResponseConfig()
	h.sendJSON(w, http.StatusOK, mock)
}

// DeleteMock handles DELETE /api/v1/mocks/:id
func (h *Handler) DeleteMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.DeleteMock(r.Context(), id); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Mock deleted successfully",
	})
}

// ActivateMock handles POST /api/v1/mocks/:id/activate
func (h *Handler) ActivateMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.ActivateMock(r.Context(), id); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Mock activated successfully",
	})
}

// DeactivateMock handles POST /api/v1/mocks/:id/deactivate
func (h *Handler) DeactivateMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.DeactivateMock(r.Context(), id); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Mock deactivated successfully",
	})
}

// TestMock handles POST /api/v1/mocks/:id/test
func (h *Handler) TestMock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get mock
	mock, err := h.store.GetMockByID(r.Context(), id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Mock not found")
		return
	}

	// Parse request
	var req models.TestMockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Convert ResponseConfig
	if err := mock.ToResponseConfig(); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate response
	startTime := time.Now()

	var responseBody interface{}
	switch mock.ResponseConfig.Type {
	case "static":
		responseBody, _ = generator.GenerateStatic(mock.ResponseConfig.Body)
	case "dynamic":
		responseBody, _ = generator.GenerateDynamic(mock.ResponseConfig.Body, req.Input)
	case "faker":
		fakerConfig := make(map[string]interface{})
		if mock.ResponseConfig.Faker != nil {
			fields := make([]interface{}, len(mock.ResponseConfig.Faker.Fields))
			for i, f := range mock.ResponseConfig.Faker.Fields {
				fields[i] = map[string]interface{}{
					"name":   f.Name,
					"type":   f.Type,
					"format": f.Format,
				}
			}
			fakerConfig["fields"] = fields
		}
		responseBody, _ = generator.GenerateFaker(fakerConfig)
	default:
		responseBody = mock.ResponseConfig.Body
	}

	tookMs := time.Since(startTime).Milliseconds()

	// Store test result
	test := &models.MockTest{
		ID:         uuid.New().String(),
		MockID:     id,
		StatusCode: mock.ResponseConfig.StatusCode,
		TookMs:     tookMs,
		CreatedAt:  time.Now(),
	}

	inputJSON, _ := json.Marshal(req.Input)
	test.Input = string(inputJSON)

	outputJSON, _ := json.Marshal(responseBody)
	test.Output = string(outputJSON)

	if err := h.store.CreateMockTest(r.Context(), test); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendJSON(w, http.StatusOK, models.TestMockResponse{
		Output:     string(outputJSON),
		StatusCode: mock.ResponseConfig.StatusCode,
		Headers:    mock.ResponseConfig.Headers,
		TookMs:     tookMs,
	})
}

// GetMockTests handles GET /api/v1/mocks/:id/tests
func (h *Handler) GetMockTests(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	tests, total, err := h.store.GetMockTestsByMockID(r.Context(), id, page, pageSize)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	h.sendJSON(w, http.StatusOK, models.ListResponse{
		Data:       tests,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// Helper methods

func (h *Handler) sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: message,
		Code:  fmt.Sprintf("%d", code),
	})
}

func (h *Handler) sendJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if err := h.store.Ping(ctx); err != nil {
		h.sendError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	h.sendJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// NotFound handles 404
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.sendError(w, http.StatusNotFound, "Not found")
}

// ServeMock handles mock requests on configured paths
func (h *Handler) ServeMock(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Get mock by path and method
	mock, err := h.store.GetMockByPathAndMethod(r.Context(), path, method)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Mock not found for this endpoint")
		return
	}

	// Convert ResponseConfig
	if err := mock.ToResponseConfig(); err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Apply delay if configured
	if mock.ResponseConfig.DelayMs > 0 {
		time.Sleep(time.Duration(mock.ResponseConfig.DelayMs) * time.Millisecond)
	}

	// Generate response
	var responseBody interface{}
	switch mock.ResponseConfig.Type {
	case "static":
		responseBody, _ = generator.GenerateStatic(mock.ResponseConfig.Body)
	case "dynamic":
		// Parse request body as input for dynamic
		input := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&input)
		responseBody, _ = generator.GenerateDynamic(mock.ResponseConfig.Body, input)
	case "faker":
		fakerConfig := make(map[string]interface{})
		if mock.ResponseConfig.Faker != nil {
			fields := make([]interface{}, len(mock.ResponseConfig.Faker.Fields))
			for i, f := range mock.ResponseConfig.Faker.Fields {
				fields[i] = map[string]interface{}{
					"name":   f.Name,
					"type":   f.Type,
					"format": f.Format,
				}
			}
			fakerConfig["fields"] = fields
		}
		responseBody, _ = generator.GenerateFaker(fakerConfig)
	default:
		responseBody = mock.ResponseConfig.Body
	}

	// Set headers
	for k, v := range mock.ResponseConfig.Headers {
		w.Header().Set(k, v)
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mock.ResponseConfig.StatusCode)
	if responseBody != nil {
		json.NewEncoder(w).Encode(responseBody)
	}
}
