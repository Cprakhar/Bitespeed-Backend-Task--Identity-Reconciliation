package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name   string
		status int
		data   interface{}
	}{
		{
			name:   "Write simple object",
			status: http.StatusOK,
			data:   map[string]string{"message": "success"},
		},
		{
			name:   "Write array",
			status: http.StatusCreated,
			data:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := WriteJSON(w, tt.status, tt.data)
			if err != nil {
				t.Errorf("WriteJSON() error = %v", err)
			}

			if w.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name        string
		jsonData    string
		expectError bool
	}{
		{
			name:        "Valid JSON",
			jsonData:    `{"name": "test", "value": 123}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonData:    `{invalid json}`,
			expectError: true,
		},
		{
			name:        "Empty JSON",
			jsonData:    `{}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")

			var result TestStruct
			err := ParseJSON(req, &result)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
