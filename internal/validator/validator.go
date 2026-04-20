package validator

import (
	"errors"
	"strings"
)

// Custom validation errors
var (
	ErrRequired  = errors.New("this field is required")
	ErrMinLength = errors.New("minimum length not met")
	ErrMaxLength = errors.New("maximum length exceeded")
	ErrInvalid   = errors.New("invalid value")
)

// Validate checks required string fields
func Validate(name, path, method string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if path == "" {
		return errors.New("path is required")
	}
	if !strings.HasPrefix(path, "/") {
		return errors.New("path must start with /")
	}
	if method == "" {
		return errors.New("method is required")
	}
	method = strings.ToUpper(method)
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	valid := false
	for _, m := range validMethods {
		if method == m {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("method must be one of: GET, POST, PUT, DELETE, PATCH")
	}
	return nil
}

// ValidateMethod validates HTTP method
func ValidateMethod(method string) error {
	method = strings.ToUpper(method)
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,
	}
	if !validMethods[method] {
		return errors.New("method must be one of: GET, POST, PUT, DELETE, PATCH")
	}
	return nil
}

// ValidateStatusCode validates HTTP status code
func ValidateStatusCode(code int) error {
	if code < 100 || code > 599 {
		return errors.New("status code must be between 100 and 599")
	}
	return nil
}

// ValidatePath validates path format
func ValidatePath(path string) error {
	if path == "" {
		return errors.New("path is required")
	}
	if !strings.HasPrefix(path, "/") {
		return errors.New("path must start with /")
	}
	return nil
}
