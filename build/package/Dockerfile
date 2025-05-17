# build/package/Dockerfile

# ---- Builder Stage ----
# Use an official Go image with Alpine Linux for a smaller base.
# Ensure the Go version matches your project's requirements.
FROM golang:1.24-alpine AS builder

# Arguments passed from cloudbuild.yaml (sourced from Taskfile's git info on submission host)
ARG BUILD_VERSION="dev-dockerfile-default"
ARG BUILD_COMMIT="unknown-dockerfile-default"
ARG MODULE_PATH="github.com/dkoosis/hello-tool-base" # Should be passed by cloudbuild.yaml
# Optional: ARG BUILD_DATE_ARG # If you choose to pass build date from cloudbuild.yaml

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker layer caching for dependencies
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download
# Verify dependencies (good practice)
RUN go mod verify

# Copy the entire application source code into the container
# This includes your ./cmd/hello-tool-base directory and other packages.
COPY . .

# Echo the received build arguments for debugging (optional, but can be helpful)
RUN echo "Dockerfile ARGs received: BUILD_VERSION='${BUILD_VERSION}', BUILD_COMMIT='${BUILD_COMMIT}', MODULE_PATH='${MODULE_PATH}'"

# Build the Go application
# The binary will be named 'hello-tool-base' and placed in the /app directory.
RUN \
    # Determine the build date: use passed ARG or generate a new one.
    # The 'date' utility is available in alpine images.
    BUILD_DATE_EFFECTIVE="${BUILD_DATE_ARG:-$(date -u '+%Y-%m-%dT%H:%M:%SZ')}" && \
    \
    # Construct LDFLAGS using the build arguments.
    # Ensure single quotes around string values in LDFLAGS if they might contain spaces,
    # though versions and commit SHAs typically don't.
    LDFLAGS_EFFECTIVE="-s -w \
    -X '${MODULE_PATH}/internal/buildinfo.Version=${BUILD_VERSION}' \
    -X '${MODULE_PATH}/internal/buildinfo.CommitHash=${BUILD_COMMIT}' \
    -X '${MODULE_PATH}/internal/buildinfo.BuildDate=${BUILD_DATE_EFFECTIVE}'" && \
    \
    echo "Dockerfile: Compiling Go application..." && \
    echo "Dockerfile: Using LDFLAGS: ${LDFLAGS_EFFECTIVE}" && \
    \
    # Execute the go build command.
    # Output to /app/hello-tool-base (WORKDIR is /app, BINARY_NAME is hello-tool-base)
    # Source path is ./cmd/hello-tool-base (relative to WORKDIR /app)
    GOOS=linux GOARCH=amd64 go build \
    -ldflags="${LDFLAGS_EFFECTIVE}" \
    -o /app/hello-tool-base \
    ./cmd/hello-tool-base

# Verification step: Check if the binary was created successfully
RUN if [ -f "/app/hello-tool-base" ]; then \
    echo "Dockerfile SUCCESS: Binary /app/hello-tool-base built successfully."; \
    else \
    echo "Dockerfile ERROR: Binary /app/hello-tool-base was NOT found after 'go build'."; \
    echo "Dockerfile: Listing contents of /app to help debug:"; \
    ls -la /app; \
    echo "Dockerfile: Listing contents of /app/cmd/hello-tool-base (source) to help debug:"; \
    ls -la /app/cmd/hello-tool-base || echo "Dockerfile: /app/cmd/hello-tool-base not found or ls failed."; \
    exit 1; \
    fi

# ---- Final Stage ----
# Use a minimal distroless image for the final application.
# base-debian12 includes CA certificates and other common libraries.
FROM gcr.io/distroless/base-debian12 AS final

# Set the working directory for the final image
WORKDIR /

# Copy only the compiled binary from the builder stage to the final image.
# This keeps the final image small and secure.
COPY --from=builder /app/hello-tool-base /app_binary

# Expose the port your application listens on (matches your Cloud Run configuration)
EXPOSE 8080

# Set default environment variables (can be overridden at runtime by Cloud Run)
ENV PORT=8080
# SERVICE_NAME is typically passed by Cloud Run if your app needs it.

# Specify the entrypoint for the container.
# Distroless images run as a non-root user 'nonroot' by default.
ENTRYPOINT ["/app_binary"]