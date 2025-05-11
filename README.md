# hello-tool-base

`hello-tool-base` is a robust Go microservice template designed for building tools intended for use with Google Vertex AI Agents. It provides a well-structured starting point with essential features for developing, testing, and deploying services on Google Cloud Run.

This template emphasizes developer experience, maintainability, and adherence to best practices for Go and cloud-native applications.

## Features

* **Basic HTTP Server:** A simple "Hello World" style HTTP server (as a starting point).
* **Structured Logging:** Configurable logging using Go's `slog` via a custom logging interface (`internal/logging`).
* **Configuration Management:** Loads configuration from YAML files and environment variables (`internal/config`).
* **Build Information:** Injects build-time variables (version, commit hash, build date) into the binary (`internal/buildinfo`).
* **Error Handling:** Defines custom application errors and maps them to appropriate responses (`internal/apperrors`), utilizing `github.com/cockroachdb/errors` for rich error context.
* **Makefile:** Comprehensive `Makefile` for common development tasks:
  * Dependency management (`deps`)
  * Code formatting (`fmt`)
  * Linting (`lint` with `golangci-lint`)
  * Testing (`test` with `gotestsum`, `test-debug`)
  * Building binaries (`build-linux`)
  * Local execution (`run-local`)
  * Docker image building and running (`docker-build-local`, `docker-run-local`)
* **Containerization:**
  * Multi-stage `Dockerfile` for optimized, secure images using a non-root distroless base.
  * `.dockerignore` for lean build contexts.
* **Cloud Deployment:**
  * `cloudbuild.yaml` for continuous integration and deployment to Google Cloud Run via Google Cloud Build.
  * Integration with Google Artifact Registry for Docker image storage.
* **OpenAPI Specification:** Basic `openapi.yaml` for defining the tool's API for Vertex AI Agent integration.

*(Upcoming/Planned Features based on best practices will include: Dedicated Health Check Endpoint, Graceful Server Shutdown, Enhanced Test Coverage, etc.)*

## Prerequisites

* [Go](https://golang.org/dl/) (version 1.22+ recommended)
* [Docker](https://www.docker.com/get-started)
* [Google Cloud SDK (gcloud CLI)](https://cloud.google.com/sdk/docs/install)
* [Make](https://www.gnu.org/software/make/)
* (Optional for some Makefile targets) `pre-commit`

## Project Structure

.
├── build/                # Docker & Cloud Build configurations
│   ├── cloudbuild/
│   │   └── cloudbuild.yaml
│   └── package/
│       └── Dockerfile
├── cmd/                  # Main application(s)
│   └── hello-tool-base/
│       └── main.go
├── configs/              # Example configuration files (if any, e.g., local.yaml)
├── docs/                 # ADRs, documentation, project context
├── internal/             # Private application logic
│   ├── apperrors/        # Application-specific error types
│   ├── buildinfo/        # Build-time information
│   ├── config/           # Configuration loading and management
│   ├── logging/          # Logging interface and setup
│   └── metrics/          # (Placeholder/Example for server metrics)
├── scripts/              # Helper scripts
├── go.mod
├── go.sum
├── LICENSE
├── Makefile              # Orchestrates development tasks
├── openapi.yaml          # API specification for Vertex AI
└── README.md

## Getting Started

1. **Clone the repository:**

    ```bash
    git clone <your-repository-url>
    cd hello-tool-base
    ```

2. **Configure Module Path:**
    * Ensure the `module` path in `go.mod` matches your intended Go module path (e.g., `github.com/your-username/hello-tool-base`).
    * Update the `MODULE_PATH` variable in the `Makefile` to match this path.

3. **Install Dependencies & Tools:**
    * Run `make deps` to tidy and download Go module dependencies.
    * Run `make install-tools` to install necessary Go development tools like `golangci-lint` and `gotestsum`.

4. **Local Development:**
    * **Run the service:** `make run-local`
        * This will build the binary and run it. The service typically listens on port 8080 by default.
        * Access in your browser or with curl: `http://localhost:8080/` or `http://localhost:8080/hello?name=YourName`
    * **Run with Docker:**
        * Build the image: `make docker-build-local`
        * Run the container: `make docker-run-local`

## Key Makefile Targets

* `make help`: Displays all available targets and their descriptions.
* `make all`: Runs checks, formats, lints, tests, and builds the binary.
* `make deps`: Tidies and downloads Go module dependencies.
* `make fmt`: Formats Go code using `golangci-lint fmt`.
* `make lint`: Runs linters using `golangci-lint`.
* `make test`: Runs unit and integration tests with `gotestsum`.
* `make test-debug`: Runs tests verbosely.
* `make build-linux`: Builds the Go binary for a Linux AMD64 environment.
* `make run-local`: Builds and runs the service locally.
* `make docker-build-local`: Builds the Docker image locally (tagged as `hello-tool-base:local`).
* `make docker-run-local`: Runs the locally built Docker image.

*(For more targets, see `make help`)*

## Configuration

The application is configured through a combination of a YAML file and environment variables.

* **Configuration File:**
  * By default, the application might look for a configuration file (e.g., `config.yaml` or specified via a flag - check `cmd/hello-tool-base/main.go`).
  * An example structure can be found in `internal/config/config.go`.
* **Environment Variables:**
  * Environment variables can override values set in the configuration file.
  * Key environment variables (refer to `internal/config/config.go` and `internal/config/applyEnvironmentOverrides`):
    * `SERVER_PORT`: Sets the port the server listens on (e.g., `8080`).
    * `SERVER_NAME`: Sets a human-readable name for the server.
    * `LOG_LEVEL`: (If implemented in logging config) Sets the logging level (e.g., `debug`, `info`, `warn`, `error`).
  * For Cloud Run deployments, environment variables (and secrets) are set via the `cloudbuild.yaml` or Cloud Run service configuration.

## Testing

* Run all tests: `make test`
* Test coverage reports are generated as `coverage.out` (and can be converted to HTML).
* Test naming conventions follow ADR-008.

## Deployment

Deployment to Google Cloud Run is automated via Google Cloud Build.

* **Trigger:** Pushes to the `main` branch (or as configured in your Cloud Build triggers and/or GitHub Actions).
* **Configuration:** Defined in `build/cloudbuild/cloudbuild.yaml`.
* **Artifacts:** Docker images are pushed to Google Artifact Registry.
* **Key Cloud Build Substitutions:**
  * `_SERVICE_NAME`: The name of the Cloud Run service (e.g., `hello-tool-base`).
  * `_GCP_REGION`: The GCP region for deployment (e.g., `us-central1`).
  * `_ARTIFACT_REGISTRY_REPO`: The name of the Artifact Registry repository.

## Contributing

*(Placeholder: Add guidelines for contributing if this project will be open to collaboration. This might include coding standards, pull request processes, etc.)*

## License

This project is licensed under the [Your License Type, e.g., Apache-2.0 License] - see the [LICENSE](LICENSE) file for details.
Next Steps for This README:

Replace placeholders:
<your-repository-url>
[Your License Type, e.g., Apache-2.0 License] - Make sure you have a LICENSE file if you specify this.
Verify Accuracy: Ensure the descriptions of features, Makefile targets, and configuration match the current state of your project.
Iterate: As you implement the other suggestions (health checks, graceful shutdown, more detailed config, etc.), update the relevant sections of this README. For example:
Add the /health endpoint to the "Features" and potentially a "Monitoring" section.
Detail new configuration options (like server timeouts).
