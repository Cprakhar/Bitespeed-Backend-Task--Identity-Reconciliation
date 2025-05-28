# Bitespeed Identity Reconciliation Service

A Go backend service for identifying and linking customer contacts across multiple purchases.

## Problem Statement

FluxKart.com needs to identify and track customer identity across multiple purchases, even when customers use different email addresses and phone numbers for each order.

## 🚀 Live Demo

**Hosted API Endpoint:** `https://your-app-name.onrender.com`

Try it out:
```bash
curl -X POST https://your-app-name.onrender.com/identify \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "phoneNumber": "1234567890"}'
```

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your settings
3. Run `make deps` to install dependencies
4. Start the server with `make run`

### Docker Development
1. Build and run with Docker Compose:
   ```bash
   make docker-run
   ```

### Deployment
1. For automatic deployment to Render.com:
   ```bash
   make deploy
   ```

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

### Docker Testing
```bash
make docker-test
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

- **Backend**: Go 1.24
- **Database**: SQLite
- **HTTP Framework**: Standard net/http
- **Deployment**: Docker, Render.com

## Development Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run unit tests
- `make test-api` - Run API tests
- `make docker-build` - Build Docker image
- `make docker-run` - Run with Docker Compose
- `make deploy` - Deploy to production
- `make help` - Show all available commands

## Project Structure

```
bitespeed-identity-reconciliation/
├── cmd/server/          # Application entry point
├── internal/
│   ├── database/        # Database connection and repository
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── pkg/utils/           # Utility functions
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── render.yaml          # Render.com deployment config
└── Makefile            # Build and deployment commands
```

## Deployment

This application is configured for easy deployment to various platforms:

- **Render.com**: Uses `render.yaml` for automatic deployment
- **Docker**: Containerized for any Docker-compatible platform
- **Local**: Simple Go binary deployment

The application automatically creates the SQLite database and required tables on startup.