# Repository Guidelines

PayIt is a Go-based Stripe integration demo. Use this guide to deliver focused contributions that keep the service reliable and secure.

## Project Structure & Module Organization

- `go.mod` defines the module `github.com/rjNemo/payit` and anchors dependencies; update it alongside any new package additions.
- Keep the HTTP server entry point in `cmd/payit/main.go`; additional CLIs belong under `cmd/<name>`.
- Place shared business logic under `internal/payments`, with Stripe client helpers in `internal/stripe`; keep packages small and unexported when possible.
- Organize web handlers, templates, and static assets under `internal/web`, `web/templates`, and `web/static`. Store sample data or fixtures beside the owning package under `testdata/`.

## Build, Test, and Development Commands

- `go run ./cmd/payit` — start the local server, loading configuration from `.env.local` when present.
- `go build ./...` — ensure every package compiles prior to opening a PR.
- `go test ./...` — execute the full unit suite; add `-run TestStripe` to focus on Stripe-specific tests during iteration.
- `STRIPE_SECRET_KEY=sk_test... go test ./internal/payments -v` — example of running package tests that require live credentials.

## Coding Style & Naming Conventions

Run `go fmt ./...` and `goimports` before committing; code should pass Go 1.25 formatting without manual tweaks. Name handlers `<Resource>Handler`, Stripe service wrappers `<Resource>Service`, and environment variables `PAYIT_*`. Centralize configuration defaults under `config/config.go` and prefer composition over inheritance-style embedding.

## Testing Guidelines

Write table-driven tests named `Test<Thing>_<Scenario>`. Co-locate doubles in `_test.go` files and keep external fixtures under `testdata/`. Target >80% coverage on payment flows, and gate integration suites behind the `-short` flag so CI can skip live Stripe runs when credentials are absent.

## Commit & Pull Request Guidelines

Use imperative, present-tense commits such as `feat: add checkout handler`; include a short body that explains why. Reference issues with `Refs #123` and note config or schema changes. PRs should summarize the change, attach manual test notes, and include screenshots or API logs when behavior shifts. Always wait for tests to pass before requesting merge.

## Security & Configuration Tips

Never commit Stripe secrets; store them in `.env.local` and rely on your process manager to inject them. Rotate exposed keys immediately. Redact cardholder data in logs and avoid persisting PII outside Stripe’s systems.
