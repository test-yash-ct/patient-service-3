# Patient Service API

Base URL: `https://<host>/v1`

## Request correlation

Clients may supply `X-Request-ID` on any request. If omitted, the service generates a UUID and returns it on the response. Structured JSON access logs include `request_id`, `tenant`, `method`, `path`, `status`, and `duration_ms`.

## Authentication

All endpoints except `/healthz`, `/meta`, and `/internal/debug/patient` (when enabled) require:

```
Authorization: Bearer <access_token>
X-Tenant-ID: <uuid>
```

## Endpoints

### `GET /meta`

Returns build and deployment metadata:

```json
{
  "service": "patient-service",
  "version": "1.2.0",
  "build_time": "2026-09-02T12:00:00Z",
  "git_sha": "abc123def"
}
```

Values are sourced from `SERVICE_VERSION`, `BUILD_TIME`, and `GIT_SHA` environment variables (defaults: `dev`, `unknown`, `unknown`).

### `GET /v1/patients/:id`

Returns a single patient record for the given identifier.

### `PATCH /v1/patients/:id`

Updates mutable fields (contact preferences, care team notes).

### `POST /v1/patients`

Creates a patient under the active tenant.

### `GET /healthz`

Liveness probe for orchestrators.

### `GET /internal/debug/patient`

Diagnostic dump of the last fetched patient payload. Availability is controlled by service configuration; see RUNBOOK.md for when this should be disabled in production.
