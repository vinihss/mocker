package validator

import (
	"testing"
)

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		method  string
		wantErr bool
	}{
		{"GET", false},
		{"POST", false},
		{"PUT", false},
		{"DELETE", false},
		{"PATCH", false},
		{"get", false},
		{"post", false},
		{"GET ", true},
		{"INVALID", true},
		{"", true},
	}

	for _, tt := range tests {
		err := ValidateMethod(tt.method)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateMethod(%q) error = %v, wantErr %v", tt.method, err, tt.wantErr)
		}
	}
}

func TestValidateStatusCode(t *testing.T) {
	tests := []struct {
		code    int
		wantErr bool
	}{
		{100, false},
		{200, false},
		{299, false},
		{400, false},
		{599, false},
		{99, true},
		{600, true},
		{0, true},
		{500, false},
	}

	for _, tt := range tests {
		err := ValidateStatusCode(tt.code)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateStatusCode(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
		}
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"/api/users", false},
		{"/api/users/123", false},
		{"/", false},
		{"api/users", true},
		{"", true},
		{"/api/v1/", false},
	}

	for _, tt := range tests {
		err := ValidatePath(tt.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		method  string
		wantErr bool
	}{
		{"", "/api", "GET", true},
		{"Test", "", "GET", true},
		{"Test", "/api", "", true},
		{"Test", "api", "GET", true},
		{"ValidName", "/api", "GET", false},
		{"ValidName", "/api", "POST", false},
		{"ValidName", "/api", "PUT", false},
		{"ValidName", "/api", "DELETE", false},
		{"ValidName", "/api", "PATCH", false},
	}

	for _, tt := range tests {
		err := Validate(tt.name, tt.path, tt.method)
		if (err != nil) != tt.wantErr {
			t.Errorf("Validate(%q, %q, %q) error = %v, wantErr %v", tt.name, tt.path, tt.method, err, tt.wantErr)
		}
	}
}
