# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set GOTOOLCHAIN to allow auto-download of newer Go versions
ENV GOTOOLCHAIN=auto

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag for Swagger documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy scripts directory for seeding data
COPY --from=builder /app/scripts ./scripts

# Copy generated Swagger documentation
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
