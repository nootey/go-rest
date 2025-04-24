# Default target runs the HTTP server
default: run

# Run the HTTP server
run:
	go run ./cmd/app server

# Build the binary
build:
	go build -o ./cmd/build/app ./cmd/app/main.go

# Run database migrations
migrate:
	go run ./cmd/app migrate $(type)

# Seed essential tables
seed:
	go run ./cmd/app seed $(type)

# Clean up binaries
clean:
	rm -f ./cmd/build/app