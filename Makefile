# Default target runs the HTTP server
default: run

# Run the HTTP server
run:
	go run cmd/app/main.go

# Build the binary (you might need to adjust based on which entrypoint you want)
build:
	go build -o cmd/build/ cmd/app/main.go

# Run database migrations
migrate:
	go run cmd/app/main.go migrate $(type)

# Seed essential tables
seed:
	go run cmd/app/main.go seed $(type)

# Clean up binaries
clean:
	rm -f ./cmd/build/app
