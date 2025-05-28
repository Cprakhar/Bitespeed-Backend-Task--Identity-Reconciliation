package services

import (
	"bitespeed-identity-reconciliation/internal/models"
    "testing"
)

func TestIdentityService_ValidateRequest(t *testing.T) {
    service := NewIdentityService()

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