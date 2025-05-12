# ---- Builder Stage ----
FROM golang:1.24-alpine AS builder

# Accept build arguments for version and commit
ARG BUILD_VERSION="dev"
ARG BUILD_COMMIT="unknown"

# Install build tools used by 'make build' if any are needed beyond Go
RUN apk add --no-cache make git

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Pass ARGs directly as make variables
RUN make build VERSION_ARG="${BUILD_VERSION}" COMMIT_ARG="${BUILD_COMMIT}"

# ---- Final Stage ----
# Use base-debian12 which includes CA certificates and other common libraries.
FROM gcr.io/distroless/base-debian12 AS final

WORKDIR /

# Copy the binary from the builder stage.
# This should match the output of 'RUN make build' in the builder stage.
COPY --from=builder /app/hello-tool-base /app_binary

EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV PORT=8080
# SERVICE_NAME will be passed by Cloud Run or Docker run command

# Distroless images run as non-root by default.
ENTRYPOINT ["/app_binary"]