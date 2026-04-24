# Architecture

## Overview

The patient service exposes a Gin HTTP server that authenticates requests, loads tenant context from headers, and reads or writes patient rows through a small repository layer on top of `pgx`.

## Components

- **cmd/server** — process entrypoint, router wiring, graceful shutdown.
- **internal/handlers** — HTTP adapters and DTO mapping.
- **internal/store** — SQL access patterns and transactions.
- **internal/auth** — bearer token parsing and subject extraction.
- **internal/middleware** — request logging, panic recovery, correlation IDs.
- **internal/config** — environment-driven configuration with conservative production defaults documented in RUNBOOK.md.

## Data flow

Clients send `Authorization: Bearer <token>` and `X-Tenant-ID`. Handlers resolve the subject identifier from the token, apply tenant scoping in queries, and return JSON payloads suitable for the scheduling and reporting services.

## Dependencies

- PostgreSQL 14+ for durable storage.
- Optional Redis (not required for baseline deployment) for future session caching.

## Failure modes

Database connectivity loss returns `503` with a stable error code. Partial writes are avoided by wrapping multi-row updates in transactions at the store layer.
