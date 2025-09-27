---
date: 2025-09-27T01:47:51+02:00
researcher: Codex
git_commit: 810a49aa4ef6e6b21cb88a6da0e8693e22cf5ae0
branch: main
repository: payit
topic: "Stripe integration demo architecture"
tags: [research, codebase, stripe, go]
status: complete
last_updated: 2025-09-27
last_updated_by: Codex
last_updated_note: Merged initial planning notes into research doc
---

# Research: Stripe integration demo architecture

## Research Question

How should we implement a Go-first web demo that lets users purchase a fake product
via Stripe in one click without adopting a frontend framework?

## Summary

A minimalist Stripe Checkout Session flow fits the project goals: serve a static
landing page from Go, expose a `/api/checkout` endpoint that creates sessions with the Stripe Go SDK, and redirect users to Stripe-hosted checkout for payment. Organize code following the guidelines in `AGENTS.md`, isolating Stripe logic under `internal/stripe`, HTTP handlers under `internal/web`, and using `cmd/payit/main.go` for bootstrapping. Local configuration should rely on environment variables (e.g., `.env.local`) to avoid embedding secrets. Automated tests can cover session creation and configuration handling with mocked Stripe clients.

## Detailed Findings

### Project Framing

- `README.md:1-8` states the app is a Stripe integration demo featuring one-time payments, so focusing on Checkout Sessions aligns with stated goals.
- `AGENTS.md:5-9` prescribes directory boundaries (`cmd/`, `internal/payments`, `internal/stripe`, `internal/web`, `web/templates`, `web/static`), providing a scaffold for new components.

### Backend HTTP Server

- Proposed entrypoint: `cmd/payit/main.go` spinning up `net/http` with routes `/` (serve landing page) and `/api/checkout` (POST to create session). Use `http.FileServer` for static assets and a custom handler for the API.
- `internal/web/server.go` (new) can construct the router, inject Stripe services, and encapsulate middleware (logging, recovery).

### Stripe Integration Layer

- `internal/stripe/client.go` (new) should wrap the official Stripe Go SDK, exposing `CreateCheckoutSession(ctx, params)` to keep handlers lightweight.
- Use environment variable `PAYIT_STRIPE_SECRET_KEY` for authentication; load via `os.LookupEnv` and fail fast if missing.

### Static Frontend Assets

- Store landing page at `web/templates/index.html` with a simple product card and a “Buy Now” button wired to POST JSON to `/api/checkout` via `fetch`.
- Serve minimal styling and JS from `web/static/` to keep the Go server responsible for asset delivery without frameworks.

### Configuration & Secrets

- Follow `AGENTS.md:17-18,26-27` guidance: run `go fmt`/`goimports`, centralize config under `config/config.go`, and never commit Stripe secrets. Support `.env.local` loading via a small helper (e.g., `github.com/joho/godotenv`) if allowed.

### Testing Strategy

- Create table-driven tests under `internal/stripe/client_test.go` to validate request payloads, stubbing Stripe client calls.
- Integration-style tests in `internal/web/server_test.go` can hit the `/api/checkout` handler using `httptest` with a fake Stripe service implementation.

## Code References

- `README.md:1-8` – Declares project purpose as a Stripe integration demo.
- `AGENTS.md:5-27` – Defines intended project layout, coding style, and security practices.

## Architecture Insights

- Adopt layered structure: handler → service (`internal/stripe`) → external Stripe API, decoupled through interfaces for testing.
- Prefer Stripe Checkout Sessions over custom payment intents to minimize client-side complexity while delivering real Stripe UX.
- Expose configuration via constructor parameters so services are testable without relying on global state.

## Planning Notes (2025-02-14)

### Objectives

- Understand user goal: Go-powered web page offering single-click Stripe checkout for a demo product without a frontend framework.
- Identify necessary backend components (HTTP server, Stripe client, config management).
- Determine minimal frontend assets required (static HTML/CSS/JS) and how to integrate Stripe Checkout or Payment Links.
- Outline testing, environment, and deployment considerations specific to this repository.

### Key Questions

1. What project structure best supports a Go-centric implementation while keeping room for future growth?
2. Which Stripe integration approach (Checkout Sessions vs. Payment Elements) balances simplicity with demo fidelity?
3. How should environment configuration and secret management be handled for local development?
4. What testing strategy ensures confidence without overcomplicating the demo?
5. Are there existing docs or notes in `./thoughts/` that can inform architectural decisions?

### Planned Research Tasks

- **Repo Survey**: Confirm current files and identify gaps for server, handlers, static assets.
- **Stripe Flow Analysis**: Compare Checkout Session vs. Payment Intent flows for a quick demo.
- **Go HTTP Server Design**: Sketch routing, handler responsibilities, and integration points.
- **Static Asset Strategy**: Determine minimal HTML/CSS/JS structure without frameworks.
- **Configuration & Secrets**: Document use of `.env` or config package for API keys.
- **Testing Considerations**: Identify unit/integration tests and Stripe mocking options.
- **Prior Research Review**: Scan `./thoughts/` for relevant history (none yet, but confirm).

### Parallel Task Outline

- Task A: Analyze current repo layout and any existing guidelines (AGENTS.md).
- Task B: Research Stripe official guidance for Go integration (if needed from prior knowledge; no external fetch unless requested).
- Task C: Design server architecture and endpoints handling Checkout session creation.
- Task D: Plan static frontend assets and their interaction with the Go backend.
- Task E: Define testing and configuration practices referencing Go ecosystem tools.

### Deliverables

- Comprehensive research document stored under `thoughts/shared/research/` following required template.
- High-level summary for user with key file/path references and recommended next steps.

## Historical Context (from ./thoughts/)

- `thoughts/2025-02-14-stripe-demo-plan.md` – Original research plan detailing objectives, tasks, and integration questions.

## Related Research

- None yet; this is the first entry in `thoughts/shared/research/`.

## Open Questions

- `hack/spec_metadata.sh` is missing, so metadata gathering requires manual commands; consider adding the script for future research compliance.
- Decide whether to embed a mock pricing catalog or keep a single hard-coded product configuration for the demo.
- Evaluate if webhooks are necessary for the demo or if redirect-based success pages suffice.
