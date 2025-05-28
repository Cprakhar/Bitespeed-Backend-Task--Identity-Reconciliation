# Bitespeed Identity Reconciliation Service

A Go backend service for identifying and linking customer contacts across multiple purchases.

## Problem Statement

FluxKart.com needs to identify and track customer identity across multiple purchases, even when customers use different email addresses and phone numbers for each order.

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your settings
3. Run `go mod tidy` to install dependencies
4. Start the server with `go run cmd/server/main.go`

## Testing

### Unit Tests
```bash
make test
```

### API Integration Tests
```bash
# Start the server in one terminal
make run

# Run API tests in another terminal
make test-api
```

### Manual Testing
```bash
# Health check
curl http://localhost:8080/health

# Create contact
curl -X POST http://localhost:8080/identify \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "phoneNumber": "1234567890"}'
```

## API Endpoints

### Health Check
```
GET /health
```
Returns `200 OK` if the service is running.

### Identity Reconciliation
```
POST /identify
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "string (optional)",
  "phoneNumber": "string (optional)"
}
```

**Response:**
```json
{
  "contact": {
    "primaryContatctId": 1,
    "emails": ["lorraine@hillvalley.edu", "mcfly@hillvalley.edu"],
    "phoneNumbers": ["123456"],
    "secondaryContactIds": [23]
  }
}
```

**Example Usage:**
```bash
curl -X POST http://localhost:8080/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "lorraine@hillvalley.edu",
    "phoneNumber": "123456"
  }'
```

## Database Schema

The service uses SQLite database with a `contacts` table for storing customer contact information.

### Contact Table Schema
- `id` - Primary key (auto-increment)
- `phone_number` - Phone number (optional)
- `email` - Email address (optional)
- `linked_id` - Foreign key to another contact (for linking)
- `link_precedence` - Either 'primary' or 'secondary'
- `created_at` - Timestamp when record was created
- `updated_at` - Timestamp when record was last updated
- `deleted_at` - Soft delete timestamp (NULL if not deleted)

## Technology Stack

- **Backend**: Go 1.21
- **Database**: SQLite
- **HTTP Framework**: Standard net/http

## Development Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run unit tests
- `make test-api` - Run API tests
- `make clean` - Clean build artifacts
- `make help` - Show all available commands