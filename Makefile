.PHONY: build run test clean help docker-build docker-run docker-test deploy

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f *.db

# Run API tests (requires server to be running)
test-api:
	./test_api.sh

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run static analysis
vet:
	go vet ./...

# Docker commands
docker-build:
	docker build -t bitespeed-identity-service .

docker-run:
	docker-compose up

docker-test:
	docker run --rm -p 8081:8080 bitespeed-identity-service &
	sleep 5
	curl -f http://localhost:8081/health
	docker stop $$(docker ps -q --filter ancestor=bitespeed-identity-service)

# Deployment
deploy:
	./deploy.sh

# Help
help:
	@echo "Available commands:"
	@echo "  build         Build the application"
	@echo "  run           Run the application"
	@echo "  test          Run unit tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  test-api      Run API integration tests"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Install dependencies"
	@echo "  fmt           Format code"
	@echo "  vet           Run static analysis"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run with Docker Compose"
	@echo "  docker-test   Test Docker container"
	@echo "  deploy        Deploy application"
	@echo "  help          Show this help message"