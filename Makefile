.PHONY: all build test clean run migrate-up migrate-down lint

# Build variables
BINARY_NAME=lozip-api
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

all: clean build test

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

test:
	go test -v -race -cover ./...

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

run:
	go run cmd/server/main.go

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

# Development tools
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker build -t lozip-api .

docker-run:
	docker run -p 8080:8080 --env-file .env lozip-api

# Test database setup
test-db-setup:
	psql -U postgres -c "DROP DATABASE IF EXISTS lozip_test;"
	psql -U postgres -c "CREATE DATABASE lozip_test;"
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/lozip_test?sslmode=disable" up

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build files"
	@echo "  make run            - Run the application"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make lint           - Run linter"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make test-db-setup  - Setup test database"
