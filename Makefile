.PHONY: run dev build test docker-up docker-down clean

# Run local development server
run:
	go run ./cmd serve-rest

# Run with Air live-reload
dev:
	air

# Build production binary
build:
	go build -ldflags="-w -s" -o bin/server ./cmd

# Run tests
test:
	go test -v ./...

# Start PostgreSQL service with Docker
docker-up:
	docker compose up -d postgres

# Stop all docker services
docker-down:
	docker compose down

# Clean temporary files
clean:
	rm -rf bin tmp
