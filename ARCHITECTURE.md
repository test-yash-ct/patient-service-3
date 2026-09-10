# Architecture

## Overview

The patient service exposes a Gin HTTP server that authenticates requests, loads tenant context from headers, and reads or writes patient rows through a small repository layer on top of `pgx`.

## Components

- **cmd/server** — process entrypoint, router wiring, graceful shutdown.
- **internal/handlers** — HTTP adapters and DTO mapping.
- **internal/store** — SQL access patterns and transactions.
- **internal/auth** — bearer token parsing and subject extraction.
- **internal/obs** — request ID middleware, structured JSON access logs, `/meta` handler.
- **internal/middleware** — thin re-export of obs access logging for legacy imports.
- **internal/config** — environment-driven configuration with conservative production defaults documented in RUNBOOK.md.

## Data flow

Clients send `Authorization: Bearer <token>` and `X-Tenant-ID`. The obs middleware assigns or propagates `X-Request-ID` and threads request ID and tenant through the request context. Handlers resolve the subject identifier from the token, apply tenant scoping in queries, and return JSON payloads suitable for the scheduling and reporting services.

## Observability

- **Request IDs** — generated when absent; echoed on responses via `X-Request-ID`.
- **Structured logs** — one JSON line per request with correlation fields.
- **Service metadata** — `/meta` exposes version and build provenance from env vars.

## Dependencies

- PostgreSQL 14+ for durable storage.
- Optional Redis (not required for baseline deployment) for future session caching.

## Failure modes

Database connectivity loss returns `503` with a stable error code. Partial writes are avoided by wrapping multi-row updates in transactions at the store layer.
