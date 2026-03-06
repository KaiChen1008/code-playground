# agent.md - code-playground (Codex)

This file defines how Codex should work in this repository.

## Project Context

This project is a self-hosted code execution platform (ideone-style) supporting Go, Python, JavaScript, Rust, and C++.

Architecture follows Clean Architecture:

- `api/`: OpenAPI/Swagger specs (source of truth).
- `cmd/server/`:
  - `delivery/`: HTTP routing, handlers, middleware (Gin).
  - `usecase/`: Application orchestration and business rules.
  - `domain/`: Core interfaces and entities.
    - `models/`: Generated from OpenAPI using `go-swagger`.
  - `repository/`: Persistence implementations (file-based in `./data`).
- `pkg/`:
  - `config/`: Viper-based config loading.
  - `errors/`: Custom error wrapping helpers.
  - `runner/`: Code execution with `testcontainers-go`.
- `ui/`: Frontend assets (Monaco editor + vanilla JS/CSS).
- `data/`: Local snippet storage.

## Codex Working Rules

- Make minimal, focused edits that match existing structure.
- Do not invent new architecture layers unless explicitly requested.
- Prefer fixing root causes over patching symptoms.
- Keep API changes spec-first:
  1. update `api/swagger.yaml`
  2. run `make swagger`
  3. update code/tests accordingly
- Never hardcode secrets, host-specific paths, or environment-specific constants.
- When touching existing code, preserve naming and package conventions.
- For non-trivial changes, run relevant checks before finishing.

## Engineering Standards

### Go

- Run `go fmt` on changed Go files.
- Keep imports grouped in this order, separated by blank lines:
  1. standard library
  2. third-party
  3. local project packages
- Run `go vet ./...` for meaningful Go changes.
- Use `pkg/errors` wrappers for contextual error handling.
- Constructor naming:
  - primary implementation: `New`
  - repository implementation: `New<Type>Repo`

### Configuration

- Add new settings in `config.yaml` and load via `pkg/config`.
- Pass only required config fields into constructors, not full config structs.

### Testing

- Add/maintain `_test.go` for repository and usecase changes.
- Use `testify/assert` or `testify/require`.
- Prefer mocks in usecase unit tests to isolate business logic.
- At minimum, run tests related to modified packages; for broad changes run full suite.

### Frontend

- Keep UI framework-free (vanilla JS/CSS + Monaco) unless explicitly requested.
- Reuse existing CSS variables/patterns.
- Preserve responsive behavior for mobile (`@media (max-width: 600px)` patterns).

## Operational Guidance

### Graceful Shutdown

When touching server lifecycle:

- use `signal.NotifyContext` for `SIGINT` and `SIGTERM`
- use `http.Server` with `BaseContext`
- keep graceful shutdown timeout (target ~5 seconds)

### Data and Execution

- Treat `data/` as persisted user content; avoid destructive changes.
- For code runner changes, maintain container isolation and timeout safeguards.

## Common Commands

- `make build`: compile server binary
- `make run`: build + run locally
- `make test`: run full test suite
- `make up`: run stack via Docker Compose
- `make swagger`: regenerate OpenAPI models

## Quality Gate Before Completion

For substantial backend changes, complete this checklist:

1. `go mod tidy`
2. `go vet ./...`
3. `go test ./...`
4. verify no secrets/local absolute paths

For scoped changes, run the smallest relevant subset and state what was run.
