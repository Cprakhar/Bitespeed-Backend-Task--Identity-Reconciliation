#!/bin/bash

echo "=== Bitespeed Identity Service Deployment ==="

# Build and test locally
echo "1. Building application..."
make build

echo "2. Running tests..."
make test

echo "3. Building Docker image..."
docker build -t bitespeed-identity-service .

echo "4. Testing Docker container..."
docker run -d --name test-container -p 8081:8080 bitespeed-identity-service
sleep 5

# Test health endpoint
if curl -f http://localhost:8081/health; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed"
    docker logs test-container
    docker stop test-container
    docker rm test-container
    exit 1
fi

# Cleanup test container
docker stop test-container
docker rm test-container

echo "5. Docker image ready for deployment!"
echo "To deploy to Render.com:"
echo "  - Push code to GitHub"
echo "  - Connect repository to Render"
echo "  - Use render.yaml for configuration"

echo "To run locally with Docker:"
echo "  docker-compose up"