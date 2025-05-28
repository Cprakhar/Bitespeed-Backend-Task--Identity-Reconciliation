package services

import (
	"bitespeed-identity-reconciliation/internal/database"
	"bitespeed-identity-reconciliation/internal/models"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestIdentityService_ValidateRequest(t *testing.T) {
	// Set up test database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer testDB.Close()

	// Create tables
	createTablesQuery := `
	CREATE TABLE contacts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phone_number TEXT,
        email TEXT,
        linked_id INTEGER,
        link_precedence TEXT NOT NULL CHECK(link_precedence IN ('primary', 'secondary')),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME,
        FOREIGN KEY (linked_id) REFERENCES contacts(id)
    );`
	if _, err := testDB.Exec(createTablesQuery); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Create service with test repository
	testRepo := database.NewContactRepository(testDB)
	service := &IdentityService{contactRepo: testRepo}

	tests := []struct {
		name        string
		request     *models.IdentifyRequest
		expectError bool
	}{
		{
			name: "Valid request with both email and phone",
			request: &models.IdentifyRequest{
				Email:       stringPtr("test@example.com"),
				PhoneNumber: stringPtr("1234567890"),
			},
			expectError: false,
		},
		{
			name: "Valid request with only email",
			request: &models.IdentifyRequest{
				Email:       stringPtr("test@example.com"),
				PhoneNumber: nil,
			},
			expectError: false,
		},
		{
			name: "Valid request with only phone",
			request: &models.IdentifyRequest{
				Email:       nil,
				PhoneNumber: stringPtr("1234567890"),
			},
			expectError: false,
		},
		{
			name: "Invalid request with empty values",
			request: &models.IdentifyRequest{
				Email:       stringPtr(""),
				PhoneNumber: stringPtr(""),
			},
			expectError: true,
		},
		{
			name: "Invalid request with nil values",
			request: &models.IdentifyRequest{
				Email:       nil,
				PhoneNumber: nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.IdentifyContact(tt.request)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
