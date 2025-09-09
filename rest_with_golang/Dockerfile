# Use official Go image as base
FROM golang:1.24.5

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose port
EXPOSE 8080

# Default command for development
CMD ["sh", "-c", "go run cmd/service/main.go"]
