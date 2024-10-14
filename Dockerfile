# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Get air
RUN go install github.com/air-verse/air@latest

# Copy go modules files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the working directory
COPY . .

# Build the Go application
RUN go build -o /app/main cmd/app/main.go

# Production Stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose the port the app will run on
EXPOSE 8080

# Run the Go binary
CMD ["/app/main"]