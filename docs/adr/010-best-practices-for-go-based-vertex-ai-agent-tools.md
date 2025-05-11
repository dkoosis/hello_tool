ADR-010: Best Practices for Go-Based Vertex AI Agent Tools

Status: Proposed
Date: 2025-05-11
1. Context

This project involves developing a suite of Go microservices ("Tools") to be invoked by a Vertex AI Agent. These tools will provide specific functionalities by interacting with enterprise systems (e.g., Namely, Stripe, Box.com) via APIs. The tools are expected to be small-to-moderate in scope and size, deployed on Google Cloud Run, and defined for the Vertex AI Agent using OpenAPI specifications. This ADR builds upon "ADR-009: Standardized Bulletproof Build Practices for Go Microservices" and adapts general Vertex AI tool best practices for our Go-centric approach. The primary goal is to ensure these tools are robust, maintainable, secure, and effectively integrated with the Vertex AI Agent.

2. Decision

We will adopt the following best practices and conventions for developing Go-based tools for the Vertex AI Agent:

2.1. Tool Design & OpenAPI Specification:

Single Responsibility: Each Go tool should ideally perform one specific, well-defined action or answer a specific type of query. This makes it easier for the Vertex AI Agent (LLM) to reason about when and how to use the tool.
Clear OpenAPI Definition (v3.x):
The openapi.yaml for each tool is the primary contract with the Vertex AI Agent.
Descriptive summary and description: Crucial for the LLM to understand the tool's purpose, capabilities, expected inputs, and outputs. Use clear, natural language. Include when the tool should and should not be used.
operationId: Must be unique and descriptive.
Parameters: Clearly define all request parameters (path, query, header, cookie) with types, formats, required status, and descriptive description fields. Provide example values.
Request Bodies: Define clearly for POST, PUT, PATCH operations.
Responses:
Define expected success responses (typically 200 OK) with a clear JSON schema for the response body. Keep response payloads concise and directly relevant to the agent's needs.
Define standard error responses (e.g., 400 Bad Request, 401 Unauthorized, 404 Not Found, 500 Internal Server Error) with a consistent JSON error structure (see 2.3).
x-google-vertex-ai-tool-name (or similar annotations): Utilize any Vertex AI-specific OpenAPI extensions for better agent integration (as seen in hello-tool-base/openapi.yaml).
Stateless Services: Go tools deployed on Cloud Run should be stateless HTTP services. Any required state should be managed externally (e.g., in a database, or passed in requests).
2.2. Go Implementation - Core Practices:

Standard Project Structure: Follow the structure outlined in "ADR-009" (e.g., cmd/, internal/handlers, internal/services, pkg/).
HTTP Handlers:
Clear separation of request parsing, business logic, and response formatting.
Use Go's standard net/http package.
Configuration:
Via environment variables (for Cloud Run).
Non-sensitive defaults can be compiled in.
Secret Management: Use Google Cloud Secret Manager for API keys and other sensitive credentials, mounted into Cloud Run as environment variables or files.
Client Libraries: Use official Go SDKs for external services (e.g., Stripe, Box) where available. Implement resilient client logic (retries, timeouts).
2.3. Error Handling (Tool to Agent):

Consistent JSON Error Responses: Tools must return structured JSON errors for client-side issues (4xx) or server-side issues (5xx) within the tool itself. A standard structure like {"error": {"code": "ERROR_CODE", "message": "Descriptive message", "details": "Optional additional details"}} is recommended. This aligns with ADR-001 principles.
Meaningful HTTP Status Codes: Use appropriate HTTP status codes (e.g., 200, 400, 401, 404, 500).
Agent Error Interpretation: The Vertex AI Agent will interpret HTTP status codes. The OpenAPI spec should clearly document what different status codes imply for the tool's execution. Errors that the agent should understand and potentially relay or act upon (e.g., "resource not found" vs. "internal tool error") should be distinguishable.
Input Validation: Rigorously validate incoming request parameters and bodies within the Go tool, returning a 400 Bad Request with clear error details if validation fails. The hello-tool-base/main.go provides a good starting point for this.
2.4. Logging & Observability:

Structured Logging: Implement structured logging (e.g., Go's log/slog) within each Go tool. Logs should be sent to Google Cloud Logging.
Key Log Information: Include request IDs (if available/traceable from the agent), relevant business identifiers, and clear error messages with stack traces (server-side only for stacks).
Cloud Run Metrics: Utilize standard Cloud Run metrics for monitoring request counts, latency, and errors.
2.5. Data Processing (e.g., gota/gota):

Default to Standard Go: For tools handling "small datasets" (e.g., a few hundred rows), default to standard Go collection types (slices, maps) for in-memory processing.
Conditional gota/gota: Evaluate gota/gota on a per-tool basis only if internal data manipulation involves multiple non-trivial, pandas-like operations (e.g., chained filters, grouping, complex sorting) that become demonstrably awkward with standard Go collections. The convenience must outweigh the added dependency.
AI Assistant Role: The AI assistant should help evaluate this trade-off, initially guiding towards standard Go solutions.
2.6. Testing:

Unit Tests: For business logic, data transformations, and helper functions. Adhere to ADR-008 for naming.
HTTP Handler Tests: Use Go's net/http/httptest for in-memory testing of HTTP handlers, verifying request parsing, response codes, and response bodies.
OpenAPI Conformance (Optional but Recommended): Consider tools to validate that actual tool responses conform to their OpenAPI schema.
Agent Simulation Testing: Use Vertex AI Agent Builder's testing/simulation features to verify the agent's interaction with the deployed tool.
2.7. Deployment & CI/CD:

Cloud Run: Deploy Go tools as services on Google Cloud Run.
Containerization: Use optimized, multi-stage Dockerfiles as per ADR-009 and hello-tool-base/Dockerfile.
CI/CD: Use Google Cloud Build (configured via cloudbuild.yaml) for automated builds, tests, and deployments, as demonstrated in hello-tool-base/cloudbuild.yaml.
3. Consequences

Positive:
Clear guidelines for developing effective and robust Go tools for Vertex AI.
Improved maintainability and consistency across different tools.
Better integration with Vertex AI Agent due to focus on OpenAPI clarity.
Enhanced security and reliability of tools.
Negative:
Requires discipline in adhering to OpenAPI best practices and structured error handling.
Initial setup for each new tool involves creating the OpenAPI spec and Go service structure, though templates can expedite this.
4. Scope

This ADR applies to all Go-based microservices developed as "Tools" to be used by the Vertex AI Agent for the Powerhouse Arts project. It complements ADR-009 by providing specific guidance for the "Tool" aspect.