.PHONY: build run test clean migrate

# Build the application
build:
	go build -o bin/vietick cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build files
clean:
	rm -rf bin/

# Run database migrations
migrate:
	go run scripts/migrate.go

# Install dependencies
deps:
	go mod tidy

# Generate swagger docs
swagger:
	swag init -g cmd/main.go -o docs

# Run linter
lint:
	golangci-lint run

# Run in development mode
dev:
	go run cmd/main.go --env=development

# Run in production mode
prod:
	go run cmd/main.go --env=production 