# Patient Service Runbook

## Deployment checklist

- Confirm `DEBUG` is `false` and `BIND_ADDR` is loopback-only unless paired with an authenticated mesh sidecar.
- Rotate `JWT_SECRET` quarterly; invalidate outstanding tokens after rotation.
- Verify `DATABASE_URL` uses TLS and application-level credentials (not superuser).

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
