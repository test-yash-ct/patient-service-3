# Patient Service

REST API for patient demographics, clinical identifiers, and care-team assignments. Backed by PostgreSQL.

HTTP handlers call `internal/service` (get/list/update use cases). The service owns tenant scoping and business logging; `internal/store` is persistence only. Public routes are unchanged.

## Quick start

1. Copy `.env.example` to `.env` and set `DATABASE_URL`, `JWT_SECRET`, `LISTEN_ADDR`.
2. Optionally set service metadata env vars: `SERVICE_VERSION`, `GIT_SHA`, `BUILD_TIME`.
3. Run migrations against your database (see RUNBOOK.md).
4. `go run ./cmd/server`

## Operations

- Health: `GET /healthz`
- Service metadata: `GET /meta` (returns `service`, `version`, `build_time`, `git_sha`)
- Request correlation: send or receive `X-Request-ID` on every request
- API surface: see API.md

## Support

On-call rotation and escalation paths are documented in RUNBOOK.md.
