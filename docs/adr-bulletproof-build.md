# ADR-009: Standardized Bulletproof Build Practices for Go Microservices

- **Status:** Proposed
- **Date:** 2025-05-10

## 1. Context

For personal and small-team Go microservice projects, there's a need for a standardized, high-quality build, test, and deployment setup that prioritizes Developer Experience (DX), reliability, security, and maintainability without introducing excessive enterprise-grade complexity. This ADR draws from the principles outlined in the "Bulletproofing Your Go Microservice Deployment (Enhanced & Integrated)" guide, tailored for smaller-scale applications where rapid iteration and ease of use are paramount. The goal is to establish a "golden path" that ensures consistency and robustness from the outset.

Existing ADRs (001-008) for the "CowGnition" project address application-specific concerns. This ADR focuses on the foundational build and deployment infrastructure applicable to new Go microservices like `hello-test-base`.

## 2. Decision

We will adopt a standardized set of build, containerization, CI/CD, and testing practices for Go microservices. These practices aim to create a "bulletproof" yet lean setup:

**2.1. Project Structure & Local Orchestration:**
    - **Opinionated Directory Layout:** Adopt a consistent project structure (e.g., `cmd/`, `internal/`, `build/package/`, `build/cloudbuild/`, `scripts/`, `configs/`).
    - **Makefile as Local Orchestration Engine:** Utilize a comprehensive `Makefile` at the project root to standardize common development tasks (lint, format, test, build binary, build Docker image, run locally, trigger deployments). This enhances DX and consistency.

**2.2. Containerization:**
    - **Optimized & Secure Dockerfiles:** Employ multi-stage Dockerfiles.
        - **Builder Stage:** Uses a Golang Alpine image for compilation and dependency resolution.
        - **Final Stage:** Uses a minimal, non-root base image (e.g., `gcr.io/distroless/static-nonroot` or `gcr.io/distroless/base-nonroot` if CGO is needed) copying only the compiled binary and necessary assets.
    - **Non-Root Execution:** Ensure containers run as non-root users.
    - **`.dockerignore`:** Maintain a comprehensive `.dockerignore` file to keep build contexts lean and secure.

**2.3. CI/CD Pipeline (Google Cloud Focused):**
    - **Cloud Build for CI/CD:** Leverage Google Cloud Build as the primary CI/CD service, configured via `cloudbuild.yaml`.
    - **GitHub Actions for Orchestration:** Use GitHub Actions workflows to trigger Cloud Build pipelines (e.g., on pushes/merges to specific branches or on tag creation).
    - **Workload Identity Federation (WIF):** Securely authenticate GitHub Actions to Google Cloud using WIF, avoiding service account key management.
    - **Dedicated Service Accounts:** Use dedicated IAM service accounts with least-privilege permissions for Cloud Build and Cloud Run services.

**2.4. Artifact Management:**
    - **Google Artifact Registry:** Store versioned Docker container images in Artifact Registry.
    - **Immutable Image Tagging:** Tag images primarily with the Git commit SHA (`$COMMIT_SHA` or `$SHORT_SHA`) for traceability and reliable rollbacks/reversions. Release tags (e.g., `v1.0.0`) can also be used. Enable tag immutability in Artifact Registry.

**2.5. Testing Integration:**
    - **Unit & Integration Tests:** Write comprehensive unit tests and relevant integration tests for the Go application.
    - **CI Test Execution:** Automatically run all tests (e.g., `go test ./...`) as a mandatory step in the Cloud Build pipeline. Fail the build if tests fail.
    - **Test Coverage (Optional but Recommended):** Generate and optionally report test coverage.
    - **Simple Post-Deployment Smoke Tests:** Include a basic smoke test in the CI/CD pipeline after deployment to Cloud Run (e.g., curling a health check or primary endpoint) to verify basic service availability.

**2.6. Configuration & Secrets Management (Basics):**
    - **Environment Variables for Configuration:** Configure Cloud Run services primarily through environment variables, adhering to Twelve-Factor App principles.
    - **Google Secret Manager for Secrets:** Store sensitive runtime secrets in Google Cloud Secret Manager and securely mount them into Cloud Run services as environment variables or files.

**2.7. Documentation:**
    - **`README.md`:** Maintain a comprehensive `README.md` detailing project setup, local development workflows (Makefile targets), build process, testing, deployment, and configuration.

## 3. Consequences

### Positive:
-   **Improved Developer Experience (DX):** Standardized tooling and workflows simplify onboarding and day-to-day development.
-   **Consistency:** Ensures all services follow the same foundational practices.
-   **Reliability:** Automated testing and deployment practices reduce manual errors and improve service stability.
-   **Security:** Secure-by-default containerization and IAM practices enhance the security posture.
-   **Maintainability:** Clear structure and automation make services easier to understand, manage, and update.
-   **Faster Onboarding:** New team members or contributors can get up to speed more quickly.
-   **Simplified Reversions:** Immutable image tagging and Cloud Run's revision management allow for straightforward rollbacks if needed.

### Negative:
-   **Initial Setup Effort:** Implementing these practices for a new project template or an existing service requires an initial time investment.
-   **Learning Curve:** Team members may need to familiarize themselves with tools like Cloud Build, Artifact Registry, and specific Makefile conventions.
-   **Tooling Overhead:** Introduces dependencies on specific cloud services (Google Cloud) and tools.
-   **Discipline Required:** Adherence to these practices requires ongoing team discipline.

## 4. Scope

This ADR primarily covers the build, test, and deployment pipeline for Go microservices. It does not dictate application-level architecture (e.g., specific Go frameworks, internal package structure beyond the top-level layout, or detailed database interaction patterns) but provides a solid foundation upon which such applications can be built. The focus is on practices that are broadly applicable and beneficial for personal to small-team sized projects hosted on Google Cloud.

## 5. Related Documents
-   "Bulletproofing Your Go Microservice Deployment (Enhanced & Integrated)" guide (Internal Document).
-   [Google Cloud Build Documentation](https://cloud.google.com/build/docs)
-   [Google Artifact Registry Documentation](https://cloud.google.com/artifact-registry/docs)
-   [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
-   [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
-   [The Twelve-Factor App](https://12factor.net/)
