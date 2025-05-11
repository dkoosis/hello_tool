# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Install make and other build essentials (like git for versioning if your Makefile uses it)
# Alpine's package manager is apk
RUN apk add --no-cache make git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod ./
COPY go.sum ./
# Ensure dependencies are downloaded. The Makefile's 'deps' target also does this,
# but doing it here explicitly can be good for layer caching.
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

WORKDIR /app # Or / for an even cleaner final stage, then CMD would be "/hello-tool-base"

# Copy the compiled binary from the builder stage.
# The Makefile's 'build' target creates './hello-tool-base' in the /app directory of the builder stage.
COPY --from=builder /app/hello-tool-base . # Copies 'hello-tool-base' to current WORKDIR (/app)

# Expose the port the application listens on (Cloud Run sets this via PORT env var)
# EXPOSE 8080 

# Command to run the executable.
CMD ["./hello-tool-base"]
