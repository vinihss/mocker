package models

import (
	"encoding/json"
	"time"
)

// Mock represents a mock endpoint configuration
type Mock struct {
	ID              string         `json:"id" db:"id"`
	Name            string         `json:"name" db:"name"`
	Path            string         `json:"path" db:"path"`
	Method          string         `json:"method" db:"method"`
	Description     string         `json:"description,omitempty" db:"description"`
	IsActive        bool           `json:"isActive" db:"is_active"`
	ResponseConfig  ResponseConfig `json:"responseConfig" db:"-"`
	ResponseType    string         `json:"-" db:"response_type"`
	ResponseStatus  int            `json:"-" db:"response_status_code"`
	ResponseHeaders string         `json:"-" db:"response_headers"`
	ResponseBody    string         `json:"-" db:"response_body"`
	ResponseDelayMs int            `json:"-" db:"response_delay_ms"`
	CreatedAt       time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time      `json:"updatedAt" db:"updated_at"`
}

// ResponseConfig defines how the mock response is generated
type ResponseConfig struct {
	Type       string            `json:"type" validate:"required,oneof=static dynamic faker"` // static, dynamic, faker
	StatusCode int               `json:"statusCode" validate:"required,min=100,max=599"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       interface{}       `json:"body" validate:"required"`
	DelayMs    int               `json:"delayMs,omitempty"`
	Faker      *FakerConfig      `json:"faker,omitempty"`
}

// FakerConfig defines faker-generated fields
type FakerConfig struct {
	Fields []FakerField `json:"fields" validate:"omitempty,dive"`
}

// FakerField defines a single faker field
type FakerField struct {
	Name   string `json:"name" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=name first_name last_name email phone phone_number uuid uuid_rfc4122 date datetime timestamp word sentence paragraph address city country state zipcode latitude longitude url ip ipv4 ipv6 credit_card credit_card_number credit_card_type"`
	Format string `json:"format,omitempty"`
}

// Instruction represents input/output transformation instructions
type Instruction struct {
	ID        string    `json:"id" db:"id"`
	MockID    string    `json:"mockId" db:"mock_id"`
	Type      string    `json:"type" db:"type"` // input, output
	Content   string    `json:"content" db:"content"`
	Order     int       `json:"order" db:"item_order"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// MockTest represents a test execution result
type MockTest struct {
	ID         string    `json:"id" db:"id"`
	MockID     string    `json:"mockId" db:"mock_id"`
	Input      string    `json:"input" db:"input"`
	Output     string    `json:"output" db:"output"`
	StatusCode int       `json:"statusCode" db:"status_code"`
	TookMs     int64     `json:"tookMs" db:"took_ms"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

// Request/Response DTOs

// CreateMockRequest defines the request for creating a mock
type CreateMockRequest struct {
	Name           string         `json:"name" validate:"required,min=1,max=100"`
	Path           string         `json:"path" validate:"required"`
	Method         string         `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Description    string         `json:"description"`
	ResponseConfig ResponseConfig `json:"responseConfig" validate:"required"`
}

// UpdateMockRequest defines the request for updating a mock
type UpdateMockRequest struct {
	Name           *string         `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Path           *string         `json:"path,omitempty" validate:"omitempty"`
	Method         *string         `json:"method,omitempty" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Description    *string         `json:"description,omitempty"`
	IsActive       *bool           `json:"isActive,omitempty"`
	ResponseConfig *ResponseConfig `json:"responseConfig,omitempty"`
}

// TestMockRequest defines the request for testing a mock
type TestMockRequest struct {
	Input map[string]interface{} `json:"input"`
}

// TestMockResponse defines the response for a mock test
type TestMockResponse struct {
	Output     string            `json:"output"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	TookMs     int64             `json:"tookMs"`
}

// ToResponseConfig converts stored DB fields to ResponseConfig
func (m *Mock) ToResponseConfig() error {
	m.ResponseConfig.Type = m.ResponseType
	m.ResponseConfig.StatusCode = m.ResponseStatus
	m.ResponseConfig.DelayMs = m.ResponseDelayMs

	if m.ResponseHeaders != "" {
		if err := json.Unmarshal([]byte(m.ResponseHeaders), &m.ResponseConfig.Headers); err != nil {
			return err
		}
	}

	if m.ResponseBody != "" {
		var bodyData map[string]interface{}
		if err := json.Unmarshal([]byte(m.ResponseBody), &bodyData); err != nil {
			return err
		}

		// Check if this is a faker config (contains "faker" key)
		if fakerData, ok := bodyData["faker"]; ok {
			m.ResponseConfig.Type = "faker"
			if fakerMap, ok := fakerData.(map[string]interface{}); ok {
				if fieldsData, ok := fakerMap["fields"]; ok {
					// Convert fields to FakerField structs
					var fields []FakerField
					if fieldsArr, ok := fieldsData.([]interface{}); ok {
						for _, f := range fieldsArr {
							if fieldMap, ok := f.(map[string]interface{}); ok {
								field := FakerField{}
								if n, ok := fieldMap["name"].(string); ok {
									field.Name = n
								}
								if t, ok := fieldMap["type"].(string); ok {
									field.Type = t
								}
								if fm, ok := fieldMap["format"].(string); ok {
									field.Format = fm
								}
								fields = append(fields, field)
							}
						}
					}
					m.ResponseConfig.Faker = &FakerConfig{Fields: fields}
				}
			}
		} else {
			m.ResponseConfig.Body = bodyData
		}
	}

	return nil
}

// FromResponseConfig converts ResponseConfig to stored DB fields
func (m *Mock) FromResponseConfig(cfg ResponseConfig) error {
	m.ResponseType = cfg.Type
	m.ResponseStatus = cfg.StatusCode
	m.ResponseDelayMs = cfg.DelayMs

	if cfg.Headers != nil {
		headers, err := json.Marshal(cfg.Headers)
		if err != nil {
			return err
		}
		m.ResponseHeaders = string(headers)
	}

	// For faker type, store the full response config including faker fields
	if cfg.Type == "faker" && cfg.Faker != nil {
		fullConfig := map[string]interface{}{
			"type":       cfg.Type,
			"statusCode": cfg.StatusCode,
			"delayMs":    cfg.DelayMs,
			"faker": map[string]interface{}{
				"fields": cfg.Faker.Fields,
			},
		}
		body, err := json.Marshal(fullConfig)
		if err != nil {
			return err
		}
		m.ResponseBody = string(body)
	} else if cfg.Body != nil {
		body, err := json.Marshal(cfg.Body)
		if err != nil {
			return err
		}
		m.ResponseBody = string(body)
	}

	return nil
}

// TransformRequest represents a transform instruction for input/output
type TransformRequest struct {
	MockID  string `json:"mockId" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=input output"`
	Content string `json:"content" validate:"required"`
	Order   int    `json:"order"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}
