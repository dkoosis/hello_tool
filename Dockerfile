# hello-tool-base/Dockerfile
# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY *.go ./

# Build the Go app
# CGO_ENABLED=0 to build a statically-linked executable
# -ldflags="-s -w" to strip debug symbols and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/hello-service .

# Start a new stage from scratch for a smaller image
FROM alpine:latest

# Add CA certificates for HTTPS calls (if your app makes any)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/hello-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/hello-service"]