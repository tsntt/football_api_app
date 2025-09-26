.PHONY: build run test clean migrate-up migrate-down docker-up docker-down

# Build the application
build-server:
	cd server && go build -o bin/server cmd/server/main.go

# Run the application
run-server:
	cd server && go run cmd/server/main.go

# Run tests
test-server:
	cd server &&  go test -v ./...

# Clean build artifacts
clean-server:
	cd server && rm -rf bin/

# Database migration up (requires golang-migrate)
migrate-up:
	echo "Make sure you have goose installed and a postgres database running on localhost:5432"
	migrate -path server/data/migrations -database "postgres://postgres:postgres@localhost:5432/football_api?sslmode=disable" up

# Database migration down
migrate-down:
	echo "Make sure you have goose installed and a postgres database running on localhost:5432"
	migrate -path server/data/migrations -database "postgres://postgres:postgres@localhost:5432/football_api?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Install dependencies
deps-server:
	cd server && go mod tidy
	cd server && go mod download

# Run with hot reload (requires air)
dev-server:
	echo "Make sure you have a postgres database running on localhost:5432"
	cd server && air

# Format code
fmt-server:
	cd server && go fmt ./...
