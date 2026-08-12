# Architecture Notes

## Failure model

The gateway assumes requests may be duplicated because of client retries, network retries, worker restarts, or upstream ambiguity.

Therefore:
- every side-effecting request must carry an idempotency key;
- terminal results are cached/persisted;
- retryability is explicit;
- deadlines are part of the request contract;
- concurrency is bounded.

## Backpressure

The demo uses a bounded semaphore. A production system should add:
- per-tenant limits;
- per-tool limits;
- queue depth metrics;
- admission control;
- overload shedding.

## Tracing

Trace identity should survive:
client -> gateway -> policy -> MCP transport -> provider.

Audit events should be append-only and independently durable.

## HashiStack

Vault fits naturally for issuing short-lived provider credentials.
Terraform owns infrastructure state.
Consul/Nomad are useful only where their operational model creates a clear advantage over native AWS primitives.
