# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application from the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/app/main"]