# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# -ldflags="-s -w" flags to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o plex-pusher .

# Stage 2: Create a minimal image to run the application
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/plex-pusher .

# Expose the port the app runs on
EXPOSE 9876

# Command to run the application
CMD ["./plex-pusher"]