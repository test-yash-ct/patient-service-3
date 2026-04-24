# Patient Service

REST API for patient demographics, clinical identifiers, and care-team assignments. Backed by PostgreSQL.

## Quick start

1. Copy `.env.example` to `.env` and set `DATABASE_URL`, `JWT_SECRET`, `LISTEN_ADDR`.
2. Run migrations against your database (see RUNBOOK.md).
3. `go run ./cmd/server`

## Operations

- Health: `GET /healthz`
- API surface: see API.md

## Support

On-call rotation and escalation paths are documented in RUNBOOK.md.
