# ---- Builder Stage ----
FROM golang:1.24-alpine AS builder

# Install build tools used by 'make build' if any are needed beyond Go
# The logs indicate 'make' and 'git' were installed.
RUN apk add --no-cache make git

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application using the Makefile's build target.
# Your Makefile's 'build' (or 'build-linux') target produces the binary '$(BINARY_NAME)',
# which is 'hello-tool-base'. This command should place it in WORKDIR (i.e., /app/hello-tool-base).
RUN make build

# ---- Final Stage ----
# Use a specific tag for distroless, for example, debian12.
# ADR-009 mentioned 'static-nonroot'.
FROM gcr.io/distroless/static-debian12 AS final
# Note: We will address the CA certificates for external HTTPS calls (loose end) later.
# If static-debian12 doesn't include them, we might switch to base-debian12 or copy them.

WORKDIR /

# Copy the binary from the builder stage.
# This should match the output of 'RUN make build' in the builder stage.
COPY --from=builder /app/hello-tool-base /app_binary

EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV PORT=8080
# SERVICE_NAME will be passed by Cloud Run or Docker run command

# Run as non-root user (distroless/*-nonroot images do this by default)
# USER nonroot:nonroot # Not needed for distroless nonroot images

ENTRYPOINT ["/app_binary"]