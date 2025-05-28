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

### 🚀 Render.com (Recommended)

The fastest way to deploy this service to production:

1. **Quick Deploy**: 
   - Update the `repo` URL in `render.yaml` with your GitHub repository
   - Push to GitHub and connect via Render Dashboard
   - Deploy automatically using the included `render.yaml` configuration

2. **Manual Setup**:
   ```bash
   # Update render.yaml with your GitHub repo URL
   # Push to GitHub
   git add .
   git commit -m "Deploy to Render"
   git push origin main
   
   # Go to render.com and create a new Blueprint service
   ```

3. **Access your deployed service**:
   - URL: `https://bitespeed-identity-service.onrender.com`
   - Health: `https://bitespeed-identity-service.onrender.com/health`

📖 **Full deployment guide**: See [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)

### 🐳 Docker

Deploy anywhere that supports Docker:

```bash
# Build and deploy locally
make docker-build
make docker-run

# Or use Docker directly
docker build -t bitespeed-identity-service .
docker run -p 8080:8080 bitespeed-identity-service
```

### 💻 Local Binary

For simple server deployment:

```bash
# Build binary
make build

# Deploy binary to server
scp bin/server user@server:/opt/bitespeed/
ssh user@server "/opt/bitespeed/server"
```

### ☁️ Other Platforms

The application works on any platform supporting:
- Docker containers (AWS ECS, Google Cloud Run, Azure Container Instances)
- Go binaries (any Linux/Windows/macOS server)
- Kubernetes (use the provided Dockerfile)

**Environment Variables for Production**:
- `PORT`: Server port (default: 8080, Render uses: 10000)
- `DB_PATH`: Database file path (default: ./contacts.db)
- `ENV`: Environment mode (production/development)

The application automatically creates the SQLite database and required tables on startup.