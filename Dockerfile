# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Install Air for live reload
RUN go install github.com/air-verse/air@latest

# Copy and install dependencies
COPY . .
RUN go mod download

# Build the app
RUN go build -o main cmd/go-rest/main.go

# Run stage
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the compiled app from builder
COPY --from=builder /app/main .
COPY --from=builder /go/bin/air /usr/bin/air
COPY . .

# Expose the port your app runs on
EXPOSE 8080

# Start Air to watch for changes
CMD ["air"]