# Makefile for HMAC Service

.PHONY: help build run test clean docker-build docker-run docker-stop deps fmt lint vet

# Default target
help:
	@echo "Available targets:"
	@echo "  help         - Show this help message"
	@echo "  deps         - Download dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run golint"
	@echo "  vet          - Run go vet"
	@echo "  test         - Run tests"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-stop  - Stop Docker container"
	@echo "  examples     - Run API examples"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run golint (install with: go install golang.org/x/lint/golint@latest)
lint:
	golint ./...

# Run go vet
vet:
	go vet ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build the application
build: deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/hmac-service ./cmd/server

# Build for current platform
build-local: deps
	go build -o bin/hmac-service ./cmd/server

# Run the application
run: build-local
	./bin/hmac-service

# Run the application with environment variables
run-dev:
	PORT=8080 LOG_LEVEL=debug go run ./cmd/server

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Docker targets
docker-build:
	docker build -t hmac-service .

docker-run: docker-build
	docker run -d --name hmac-service -p 8080:8080 hmac-service

docker-stop:
	docker stop hmac-service || true
	docker rm hmac-service || true

docker-logs:
	docker logs -f hmac-service

# Docker Compose targets
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

compose-logs:
	docker-compose logs -f

# Run API examples (requires the service to be running)
examples:
	chmod +x examples/curl_examples.sh
	./examples/curl_examples.sh

# Development workflow
dev-setup: deps fmt vet test

# CI/CD workflow
ci: deps fmt vet test build

# Install development tools
install-tools:
	go install golang.org/x/lint/golint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Security check (requires gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
security:
	gosec ./...