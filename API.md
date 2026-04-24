# Patient Service API

Base URL: `https://<host>/v1`

## Authentication

All endpoints except `/healthz` and `/internal/debug/patient` (when enabled) require:

```
Authorization: Bearer <access_token>
X-Tenant-ID: <uuid>
```

## Endpoints

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
