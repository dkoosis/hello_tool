Project Context & Guiding Principles: Slack Enterprise Chatbot for Powerhouse Arts (Go/Vertex AI)

AI Assistant: Gemini
Date Context Last Updated: May 11, 2025
Version: 2.0

1. Overall Goal & Core User Need

Develop an enterprise chatbot integrated with Slack for Powerhouse Arts. The chatbot will provide timely, accurate, and domain-relevant answers to employee questions by securely accessing and processing data from company systems (Namely, Stripe, Box.com documents).

Target Scale: Small enterprise, serving a group of a few hundred internal users.
Core Value: Enable users to ask questions naturally and get contextualized answers via an LLM-powered agent.
2. Core Trade-off Philosophy

We consistently prioritize Developer Experience (DX), iteration speed, and manageable complexity over achieving extreme enterprise-grade robustness (e.g., 100% uptime, complex rollback/disaster recovery, massive scalability) for initial phases. When faced with a choice, opt for the simpler, more maintainable solution suitable for our scale and DX focus.

AI Role: Actively help identify and suggest these simpler paths. Challenge suggestions that might lead to premature complexity for our defined scope.
3. Guiding Principles for Development & AI Collaboration

Simplicity & Iteration:
Preference: Start with the simplest viable MVP. Add complexity only when justified.
AI Assistance: Help identify the most straightforward path. If a complex solution is suggested, please also offer a simpler alternative.
Robustness (for Core Functionality):
Focus: Clear error handling in Go tools (structured JSON errors, meaningful HTTP codes), input validation against OpenAPI specs, and secure secret management (Google Cloud Secret Manager).
AI Assistance: Provide Go patterns for these, referencing relevant ADRs.
Maintainability & Code Quality:
Focus: Logical modularity in Go tools, unit tests for key logic (ADR-008 naming, httptest), and structured logging (log/slog).
AI Assistance: Suggest idiomatic Go, help refactor, provide testing/logging examples.
Pragmatism over Premature Optimization/Complexity:
Data Manipulation (Small Datasets): Default to standard Go collections. Evaluate gota/gota only if internal tool logic becomes demonstrably complex for pandas-like operations within a single tool.
AI Assistance (Gota): Remind me of this default. Help evaluate if gota/gota is justified for a specific tool's small dataset, or if it's premature optimization.
State Management: Go tools should be stateless HTTP services.
Advanced Agent Concepts: Defer deep dives into ADK, Agentspace, or complex LLM-generated SQL until core functionality is proven and explicitly justified.
AI Assistance: Help evaluate when a more complex pattern is truly needed versus keeping it simple.
4. AI Assistant Collaboration Model

Role: Act as a coding partner, sounding board for design decisions, provider of Go & Google Cloud best practices, and debugger.
Interaction: I will provide context (like this document), ADRs, specific tasks, and code. Please offer code generation, examples, explanations, critiques, and help align with these guiding principles.
Primary Focus: Help build individual Go tools (including their openapi.yaml) effectively for integration with Vertex AI Agent Builder.
5. Core Technologies (Summary)

Conversational AI & Orchestration: Vertex AI Agent Builder
Document Search: Vertex AI Search (for Box.com)
Custom Data Access/Actions: Go microservices ("Tools") on Google Cloud Run
Tool API Definition: OpenAPI Specification (YAML v3.x)
Slack Integration: Slack Bolt SDK (likely Python/Node.js on Cloud Run - separate from Go tools)
CI/CD & Deployment: Docker, Google Cloud Build, Artifact Registry
Secret Management: Google Cloud Secret Manager
6. Relevant ADRs (Key Guiding Documents)

ADR-001: Error Handling Strategy
ADR-002: Schema Validation Strategy (for tools via OpenAPI)
ADR-005: Secret Management (adapted for Cloud Run)
ADR-006: Modular Multi-Service Support (informing Go tool structure)
ADR-008: Test Naming Convention
ADR-009: Standardized Bulletproof Build Practices
ADR-010: Best Practices for Go-Based Vertex AI Agent Tools