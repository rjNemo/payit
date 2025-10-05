# Stripe Checkout Demo Implementation Plan

## Overview

Implement a Go-based Stripe Checkout demo that serves a static landing page for a single hardcoded product and uses Stripe-hosted checkout to process payments, aligning with the previously captured research.

## Current State Analysis

- Repository currently contains only `go.mod`, `README.md`, and `AGENTS.md`; there is no executable application code or directory structure yet (`README.md:1-8`, `AGENTS.md:5-27`).
- No configuration handling, HTTP server, or Stripe integration exists; we must introduce all components from scratch.
- Research document `thoughts/shared/research/2025-09-27-stripe-integration-demo.md` outlines recommended architecture, folder layout, and open questions.

## Desired End State

Deliver a runnable Go service (`cmd/payit/main.go`) that serves a static product page, exposes a `/api/checkout` endpoint to create a Stripe Checkout Session for a single hardcoded product, and documents setup/testing steps. The codebase should follow the structure in `AGENTS.md`, with automated tests covering the Stripe service wrapper and HTTP handler behavior.

### Key Discoveries

- `thoughts/shared/research/2025-09-27-stripe-integration-demo.md:21-55` – Defines target architecture for server, Stripe layer, assets, and testing.
- `AGENTS.md:5-18` – Specifies preferred directory layout (`cmd/`, `internal/`, `web/`, `config/`) and formatting expectations.
- `AGENTS.md:26-27` – Emphasizes proper handling of Stripe secrets via environment variables.

## Out of Scope

- Supporting multiple or dynamic products; only a single hardcoded demo product is required.
- Implementing Stripe webhooks or fulfillment workflows; success relies on Stripe redirect pages.
- Adding auxiliary tooling such as `hack/spec_metadata.sh` or deployment scripts.

## Implementation Approach

Build the application in four incremental phases: establish project skeleton and configuration, implement Stripe Checkout backend logic, wire static frontend assets with routing, and finish with testing plus documentation updates. Each phase will introduce isolated packages, enabling focused testing and straightforward iteration.

## Phase 1: Project Skeleton & Configuration

### Overview

Create the application structure, entrypoint, and configuration loader for environment variables.

### Changes Required

**File**: `go.mod`  
**Changes**: Add required module dependencies (`github.com/stripe/stripe-go/v83`, optional `github.com/joho/godotenv` if used) and tidy module.

**File**: `cmd/payit/main.go` (new)  
**Changes**: Bootstrap configuration, initialize logger, construct HTTP server (delegated to `internal/web`).

**File**: `config/config.go` (new)  
**Changes**: Define `Config` struct holding Stripe secret key, publishable key, and product metadata (name, price, currency, success/cancel URLs). Load values from environment (with optional `.env.local` support) and validate presence.

**File**: `config/config_test.go` (new)  
**Changes**: Unit tests ensuring configuration loading validates required fields and default behaviors.

**File**: Directory scaffolding (`internal/stripe/`, `internal/web/`, `web/templates/`, `web/static/`, `testdata/`)  
**Changes**: Create empty placeholder files or README stubs as needed so later phases can populate them.

```go
// config/config.go
package config

type ProductConfig struct {
    Name        string
    Description string
    PriceCents  int64
    Currency    string
    SuccessURL  string
    CancelURL   string
}

type Config struct {
    StripeSecretKey     string
    StripePublishableKey string
    Product             ProductConfig
}
```

### Success Criteria

#### Automated Verification

- [x] `go fmt ./...`
- [x] `go build ./...`

#### Manual Verification

- [x] Application starts with placeholder server: `go run ./cmd/payit` _(bind restricted in sandbox; confirmed startup log before failure)_
- [x] Missing environment variables cause a clear startup error

## Phase 2: Stripe Checkout Backend

### Overview

Implement the Stripe client wrapper, checkout session creator, and API handler.

### Changes Required

**File**: `internal/stripe/client.go` (new)  
**Changes**: Wrap Stripe SDK, expose interface for session creation, configure Stripe API key.

**File**: `internal/stripe/types.go` (new)  
**Changes**: Define request/response structs (e.g., `CheckoutSessionRequest`, `CheckoutSessionResult`).

**File**: `internal/stripe/client_test.go` (new)  
**Changes**: Table-driven tests using a fake Stripe API client to validate payload construction.

**File**: `internal/web/handlers.go` (new)  
**Changes**: Implement `/api/checkout` HTTP handler accepting POST, calling Stripe service, returning JSON (session ID / URL) with proper error handling.

**File**: `internal/web/server.go` (new)  
**Changes**: Build `http.Handler` wiring API routes and injecting Stripe service.

**File**: `internal/web/handlers_test.go` (new)  
**Changes**: Use `httptest` with a mocked Stripe service to assert status codes, response payloads, and error cases.

```go
// internal/web/handlers.go
func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    session, err := h.CheckoutService.CreateSession(ctx)
    if err != nil {
        http.Error(w, "checkout session failed", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(session)
}
```

### Success Criteria

#### Automated Verification

- [ ] `go test ./internal/stripe`
- [ ] `go test ./internal/web`

#### Manual Verification

- [ ] `curl -X POST http://localhost:8080/api/checkout` returns session payload with mock keys
- [ ] Error cases logged and surfaced with 500 response when Stripe call fails

## Phase 3: Static Frontend & Routing

### Overview

Implement landing page, static assets, and route wiring for root and assets.

### Changes Required

**File**: `web/templates/index.html` (new)  
**Changes**: HTML page describing the product with “Buy Now” button invoking JavaScript to POST to `/api/checkout` and redirect to returned `url`.

**File**: `web/static/styles.css` (new)  
**Changes**: Minimal styling for the product page.

**File**: `web/static/app.js` (new)  
**Changes**: JS `fetch` call to `/api/checkout`, handles response, redirects or displays error.

**File**: `internal/web/server.go`  
**Changes**: Serve template on `/`, static assets via `http.FileServer`, inject publishable key into template data if needed.

**File**: `internal/web/templates.go` (new)  
**Changes**: Helper to parse templates and render with provided data (publishable key, product info).

**File**: `internal/web/server_test.go` (new)  
**Changes**: Verify root route renders successfully and static assets are served.

### Success Criteria

#### Automated Verification

- [ ] `go test ./internal/web`

#### Manual Verification

- [ ] Visit `http://localhost:8080/` and confirm product page renders with publishable key available to JS
- [ ] Clicking “Buy Now” redirects to Stripe Checkout when using valid keys

## Phase 4: Testing & Developer Experience

### Overview

Finalize tests, documentation, and developer onboarding.

### Changes Required

**File**: `README.md`  
**Changes**: Add setup instructions (env vars, running server, creating test mode keys), testing commands, and manual verification steps.

**File**: `.env.example` (new)  
**Changes**: Document expected environment variables (`PAYIT_STRIPE_SECRET_KEY`, `PAYIT_STRIPE_PUBLISHABLE_KEY`, `PAYIT_PRODUCT_NAME`, etc.).

**File**: `Makefile` (new, optional)  
**Changes**: Provide convenience targets (`make run`, `make test`) if deemed helpful for onboarding.

**File**: `internal/stripe/mock.go` (new)  
**Changes**: Provide simple mock implementation used in tests to avoid Stripe network calls.

**File**: `thoughts/shared/plans/2025-09-27-stripe-checkout-demo.md`  
**Changes**: Update checkboxes for completed phases as work progresses.

### Success Criteria

#### Automated Verification

- [x] `go fmt ./...`
- [ ] `go test ./...`
- [x] `go build ./...`

#### Manual Verification

- [ ] README instructions reproduce successful checkout flow in Stripe test mode
- [ ] `.env.example` allows quick setup with test keys

## Testing Strategy

- Unit tests for configuration loader, Stripe client wrapper, and HTTP handlers (`config`, `internal/stripe`, `internal/web`).
- Integration-style handler tests using mocks to simulate Stripe responses.
- Manual validation of end-to-end flow with Stripe test keys in browser.

## Performance Considerations

- Single product flow with minimal traffic; standard `net/http` defaults suffice. Ensure Stripe client is reused to avoid unnecessary overhead.

## Migration Notes

- No existing users; deployment is greenfield. Ensure new directories and files are committed together.

## References

- Research: `thoughts/shared/research/2025-09-27-stripe-integration-demo.md`
- Guidelines: `AGENTS.md`
- Ticket (README context): `README.md`
