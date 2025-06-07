.PHONY: run run-dev run-prod test test-dev test-prod build docker-build

# Development environment
run-dev:
	ENVIRONMENT=development go run cmd/api/main.go

# Production environment
run-prod:
	ENVIRONMENT=production go run cmd/api/main.go

# Default to development
run: run-dev

# Test commands
test-dev:
	ENVIRONMENT=development go test -v ./...

test-prod:
	ENVIRONMENT=production go test -v ./...

# Default test uses test environment
test:
	ENVIRONMENT=test go test -v ./...

build:
	go build -o bin/api cmd/api/main.go

docker-build:
	docker build -t lean-backend-boilerplate-golang .
