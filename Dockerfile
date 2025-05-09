# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod ./
COPY go.sum ./
# Ensure dependencies are downloaded. The Makefile's 'deps' target also does this.
# Running 'go mod download' here is belt-and-suspenders, or Makefile's deps can be used.
RUN go mod download

# Copy the entire application source code, including the Makefile
COPY . .

# Build the application using the Makefile.
# This will create an executable named 'app' (or BINARY_NAME from Makefile) in the /app directory.
RUN make build

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS calls if your app makes them
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary from the builder stage.
# The Makefile's 'build' target creates './app' in the /app directory of the builder stage.
COPY --from=builder /app/app .

# Expose the port the application listens on (Cloud Run sets this via PORT env var)
# EXPOSE 8080 # This is more for documentation; Cloud Run handles port exposure.

# Command to run the executable.
# Cloud Run will inject the PORT environment variable.
CMD ["./app"]