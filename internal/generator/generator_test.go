package generator

import (
	"testing"
)

func TestGenerateStatic(t *testing.T) {
	body := map[string]interface{}{
		"message": "Hello World",
		"status":  "success",
	}

	result, err := GenerateStatic(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	if resultMap["message"] != "Hello World" {
		t.Errorf("expected message 'Hello World', got %v", resultMap["message"])
	}
}

func TestGenerateDynamic(t *testing.T) {
	// Template with input placeholder
	template := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "{{input.name}}",
			"email": "{{input.email}}",
		},
		"timestamp": "{{timestamp}}",
	}

	input := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	result, err := GenerateDynamic(template, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	userMap, ok := resultMap["user"].(map[string]interface{})
	if !ok {
		t.Fatal("expected user map")
	}

	if userMap["name"] != "John Doe" {
		t.Errorf("expected name 'John Doe', got %v", userMap["name"])
	}
}

func TestGenerateDynamicWithTimestamps(t *testing.T) {
	template := map[string]interface{}{
		"date":      "{{date}}",
		"datetime":  "{{datetime}}",
		"timestamp": "{{timestamp}}",
		"uuid":      "{{uuid}}",
	}

	input := map[string]interface{}{}

	result, err := GenerateDynamic(template, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	// Check that date, datetime, timestamp, uuid are present
	if resultMap["date"] == nil {
		t.Error("expected date to be present")
	}
	if resultMap["datetime"] == nil {
		t.Error("expected datetime to be present")
	}
	if resultMap["timestamp"] == nil {
		t.Error("expected timestamp to be present")
	}
	if resultMap["uuid"] == nil {
		t.Error("expected uuid to be present")
	}
}

func TestGenerateFaker(t *testing.T) {
	fakerConfig := map[string]interface{}{
		"fields": []interface{}{
			map[string]interface{}{
				"name": "name",
				"type": "name",
			},
			map[string]interface{}{
				"name": "email",
				"type": "email",
			},
			map[string]interface{}{
				"name": "phone",
				"type": "phone",
			},
		},
	}

	result, err := GenerateFaker(fakerConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	// Check that all fields are generated
	if resultMap["name"] == nil {
		t.Error("expected name field")
	}
	if resultMap["email"] == nil {
		t.Error("expected email field")
	}
	if resultMap["phone"] == nil {
		t.Error("expected phone field")
	}
}

func TestGenerateFakerFields(t *testing.T) {
	testCases := []struct {
		fieldType string
	}{
		{"name"},
		{"first_name"},
		{"last_name"},
		{"email"},
		{"phone"},
		{"uuid"},
		{"date"},
		{"datetime"},
		{"word"},
		{"sentence"},
		{"paragraph"},
		{"city"},
		{"country"},
		{"state"},
		{"zipcode"},
		{"latitude"},
		{"longitude"},
		{"url"},
		{"ipv4"},
		{"ipv6"},
		{"credit_card"},
		{"credit_card_type"},
	}

	for _, tc := range testCases {
		fakerConfig := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": tc.fieldType,
					"type": tc.fieldType,
				},
			},
		}

		result, err := GenerateFaker(fakerConfig)
		if err != nil {
			t.Errorf("unexpected error for type %s: %v", tc.fieldType, err)
			continue
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Errorf("expected map result for type %s", tc.fieldType)
			continue
		}

		if resultMap[tc.fieldType] == nil {
			t.Errorf("expected field %s to be generated", tc.fieldType)
		}
	}
}

func TestGenerateDynamicPreservesStructure(t *testing.T) {
	// Test that nested structures are preserved
	template := map[string]interface{}{
		"data": map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"id": 1, "name": "{{input.item1}}"},
				map[string]interface{}{"id": 2, "name": "{{input.item2}}"},
			},
		},
	}

	input := map[string]interface{}{
		"item1": "First Item",
		"item2": "Second Item",
	}

	result, err := GenerateDynamic(template, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	// Check that result has data key
	if resultMap["data"] == nil {
		t.Error("expected data to be present")
	}
}

func BenchmarkGenerateStatic(b *testing.B) {
	body := map[string]interface{}{
		"message": "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateStatic(body)
	}
}

func BenchmarkGenerateFaker(b *testing.B) {
	fakerConfig := map[string]interface{}{
		"fields": []interface{}{
			map[string]interface{}{"name": "name", "type": "name"},
			map[string]interface{}{"name": "email", "type": "email"},
			map[string]interface{}{"name": "phone", "type": "phone"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateFaker(fakerConfig)
	}
}

func BenchmarkGenerateDynamic(b *testing.B) {
	template := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "{{input.name}}",
		},
	}
	input := map[string]interface{}{
		"name": "John",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateDynamic(template, input)
	}
}
