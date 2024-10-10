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

# Expose the application port
EXPOSE 3000

# Start the application
CMD [ "air", "-c", ".air.toml" ]