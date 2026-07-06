# CLAUDE.md — SAME Message to Mesh

Ground rules for all development (human and AI) in this repository.

## Stack

- **Backend**: Go, Postgres, `goose`
- **Frontend**: TypeScript, React, Vite
- **Infrastructure**: Docker / docker-compose
- **AI files**: `.ai/` directory (prompts, specs, decisions, plan)

## Architecture — Hexagonal (Ports & Adapters)

All backend code follows hexagonal architecture. Dependency direction is always inward:

```
adapter → application → domain
```

### Layer rules

| Layer | Path | Rule |
|---|---|---|
| Domain | `backend/internal/domain/` | Pure Go only. Zero framework or infrastructure imports. Contains entities and port interfaces. |
| Application | `backend/internal/application/` | Orchestrates domain via port interfaces. No HTTP, DB, or external service imports. |
| Adapter | `backend/internal/adapter/` | All framework/infra code. Implements or consumes port interfaces. |

The frontend mirrors this structure under `frontend/src/` with `domain/`, `adapters/`, and `ui/` layers.

## Configuration — 12-Factor

- All runtime config comes from environment variables. No hardcoded values.
- `.env.example` is the canonical reference for required variables. Commit changes to it.
- `.env` is gitignored. Never commit it.
- Local dev loads env vars via docker-compose `env_file`.

## Database

- Postgres only.
- Migrations managed with `goose` using plain SQL files with `-- +goose Up` / `-- +goose Down` annotations (`backend/internal/adapter/repository/migrations/`).
- **Migrations run automatically at application startup** before the HTTP server starts. `main.go` calls `goose.Up()` on boot.
- New migrations: `make migrate-create NAME=<description>`.
- **All application state must be persisted in Postgres.** No in-memory-only state. Container restarts must be fully transparent.

## MQTT

- Controlled by `MQTT_ENABLED` env var (default `false`).
- When `MQTT_ENABLED=false`, the MQTT publisher must not start and the application must boot normally without errors.
- Phase 1 does not require MQTT. Phase 2 adds publishing.

## SDR

- The SDR adapter wraps `rtl_fm | multimon-ng` via `exec.Command`.
- The SDR device path and frequency are persisted in the DB (sdr_config table) and configurable via UI.
- SDR adapter implements a port interface to allow mocking in unit tests.

## Logging

- Use `log/slog` (stdlib). No third-party logging libraries.
- Always use the JSON handler. Log level set via `LOG_LEVEL` env var (`debug`, `info`, `warn`, `error`). Default: `info`.
- Never log sensitive data (passwords, tokens, PII).

## Testing

- Every package must have tests.
- Backend target: **>90% coverage**. `testing` stdlib + `testify/assert` + `testify/require` + `testify/mock`. Use `mockery` to generate mocks from port interfaces.
- Frontend target: **>80% coverage**. Vitest + React Testing Library + MSW for API mocking.
- Domain and application layers must be unit tested in isolation — no DB, no HTTP.
- Run all tests: `make test`. Run with coverage: `make coverage`.

## Pre-Commit Requirements

Both must pass before any commit:

```
make lint
make fmt
```

- Backend linter: `golangci-lint`
- Frontend linter: ESLint
- Backend formatter: `gofmt`
- Frontend formatter: Prettier

## Git Workflow

- NEVER commit directly to `main`/`master`. Always create a new branch first.
- Before making changes, check the current branch with `git branch --show-current`. If on `main`/`master`, create and switch to a new branch before editing any files.
- Branch naming convention: `feature/<short-description>`, `fix/<short-description>`, or `chore/<short-description>`.
- Do not create a new branch if one was already created for this task/session — reuse it.

## Plan File

`.ai/PLAN.md` is the source of truth for project progress. Keep it updated:

- Check off completed work as it is done.
- Log architecture decisions with rationale.
- Add new backlog items as they are identified.
