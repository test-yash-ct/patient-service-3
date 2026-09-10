# Patient Service Runbook

## Deployment checklist

- Confirm `DEBUG` is `false` and `BIND_ADDR` is loopback-only unless paired with an authenticated mesh sidecar.
- Set `SERVICE_VERSION`, `GIT_SHA`, and `BUILD_TIME` from the CI pipeline so `/meta` reflects the running build.
- Rotate `JWT_SECRET` quarterly; invalidate outstanding tokens after rotation.
- Verify `DATABASE_URL` uses TLS and application-level credentials (not superuser).

## Environment variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `SERVICE_VERSION` | Semantic version exposed on `/meta` | `dev` |
| `GIT_SHA` | Commit SHA exposed on `/meta` | `unknown` |
| `BUILD_TIME` | ISO-8601 build timestamp on `/meta` | `unknown` |

## Request tracing

- Ingress and mesh sidecars should forward `X-Request-ID` when present.
- Search logs by `request_id` in the observability stack to correlate handler, auth, and store events.

## Common incidents

### Elevated 5xx from database

1. Check connection pool saturation in metrics.
2. Fail over read traffic to replica if configured.
3. Scale replicas horizontally if CPU-bound.

### Slow patient lookups

1. Validate indexes on `(tenant_id, patient_id)`.
2. Review recent query plans in the observability stack.

## Rollback

Re-deploy the previous container image tag recorded in the change ticket. Runbook owner: platform on-call.
