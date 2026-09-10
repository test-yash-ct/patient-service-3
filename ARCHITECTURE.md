# Architecture

## Overview

The patient service exposes a Gin HTTP server that authenticates requests and delegates tenant-scoped patient use cases to an in-process domain service layer. Persistence stays in the store package on top of `pgx`.

## Components

- **cmd/server** — process entrypoint, constructor wiring, graceful shutdown.
- **internal/handlers** — HTTP adapters: bind headers/path, call the service, map errors to status codes.
- **internal/service** — domain use cases (`Patients`: get, list, update). Tenant resolution, request-ID business logging, and access rules live here.
- **internal/store** — SQL access patterns and transactions. Queries are always scoped by `tenant_id`.
- **internal/auth** — bearer token parsing and subject extraction.
- **internal/obs** — request ID middleware, structured JSON access logs, `/meta` handler.
- **internal/middleware** — thin re-export of obs access logging for legacy imports.
- **internal/config** — environment-driven configuration with conservative production defaults documented in RUNBOOK.md.

## Data flow

Clients send `Authorization: Bearer <token>` and `X-Tenant-ID`. The obs middleware assigns or propagates `X-Request-ID` and threads request ID (and header tenant when present) through the request context.

1. Handler parses the bearer token (transport concern).
2. Domain service resolves tenant from claims vs header (mismatch is forbidden), applies business rules, and logs access with `request_id` (no PHI in logs).
3. Store executes tenant-scoped SQL and returns rows.
4. Handler maps typed service errors to HTTP status codes and returns JSON.

This layering is in-process. It is not a new microservice.

## Observability

- **Request IDs** — generated when absent; echoed on responses via `X-Request-ID`.
- **Structured logs** — one JSON line per request at the obs layer; domain events logged from `internal/service`.
- **Service metadata** — `/meta` exposes version and build provenance from env vars.

## Dependencies

- PostgreSQL 14+ for durable storage.
- Optional Redis (not required for baseline deployment) for future session caching.

## Failure modes

Database connectivity loss returns `503` with a stable error code. Partial writes are avoided by wrapping multi-row updates in transactions at the store layer.
