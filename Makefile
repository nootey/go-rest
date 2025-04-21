# Default target runs the HTTP server
default: run

# Run the combined server (both HTTP and gRPC)
run:
	go run cmd/app/main.go

# Build the binary (you might need to adjust based on which entrypoint you want)
build:
	go build -o cmd/build/ cmd/app/main.go
