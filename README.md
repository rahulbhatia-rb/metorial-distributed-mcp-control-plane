# Metorial Distributed MCP Control Plane

A role-specific proof-of-concept inspired by Metorial's public product surface and the distributed-systems role description.

## What this demonstrates

Metorial connects employees and AI agents to tools while enforcing governance, access control, and tracing. This project explores the infrastructure problem underneath that surface: **how a distributed tool-execution layer can stay safe, observable, idempotent, and resilient under retries and partial failure.**

The project is intentionally not a clone of Metorial and makes no assumptions about Metorial's private architecture.

## Architecture

```text
AI client / agent
      |
      v
TypeScript control-plane client
      |
      v
Go execution gateway
  |       |       |
policy  tracing  idempotency
  |       |       |
  +---- provider/tool adapter ----> MCP/tool backend
```

### Key properties

- Per-request actor + agent identity
- Capability-based policy checks
- Idempotency keys for side-effecting tool calls
- Bounded concurrency / backpressure
- Retry classification
- Deadline propagation
- Trace/span propagation
- Structured audit events
- Deterministic failure states
- TypeScript client for control-plane integration

## Why it maps to the role

The role asks for:
- deep distributed-systems understanding
- TypeScript and Go
- AWS + HashiStack / advanced DevOps
- large production systems
- AI agents, MCP, and tool calls

This repo focuses directly on those boundaries rather than presenting a generic CRUD demo.

## Go execution model

A request includes:
- `actor_id`
- `agent_id`
- `tool`
- `operation`
- `arguments`
- `idempotency_key`
- `deadline_ms`
- `trace_id`

Before execution the gateway:
1. validates the request;
2. evaluates capability policy;
3. checks idempotency state;
4. acquires bounded execution capacity;
5. creates an audit/trace event;
6. executes the tool adapter;
7. classifies failures as retryable or terminal;
8. records the terminal result against the idempotency key.

## TypeScript client

`ts-client/` shows how an application or agent runtime could call the gateway while preserving:
- trace IDs
- actor identity
- agent identity
- idempotency keys
- explicit deadlines

## Example

```bash
go test ./...
go run ./cmd/gateway
```

In a second terminal:

```bash
cd ts-client
npm install
npm run demo
```

## Production evolution

For a production-grade version I would add:
- Envoy/NLB ingress with mTLS
- AWS ECS/EKS deployment
- service discovery
- Redis or DynamoDB-backed idempotency store
- SQS/Kafka for deferred execution where appropriate
- OpenTelemetry collector + trace sampling policies
- Vault for dynamic credentials
- Consul for discovery/config where justified
- Nomad as an alternate HashiStack scheduling path
- per-tenant concurrency budgets
- circuit breakers per provider/tool
- adaptive retry budgets
- durable outbox for audit events
- regional failover
- chaos testing
- replay tooling for failed tool calls
- schema/version compatibility checks for MCP tools

## AWS + HashiStack reference path

A pragmatic deployment could use:
- AWS NLB/ALB
- ECS/EKS or Nomad
- RDS/DynamoDB for control-plane state
- ElastiCache/Redis for hot idempotency state
- S3 for long-term trace/audit retention
- CloudWatch + OpenTelemetry
- Vault for secrets and dynamic credentials
- Terraform for infrastructure provisioning
- Consul only where service discovery/config requirements justify it

## Design tradeoffs

### At-least-once vs exactly-once

Exactly-once execution across distributed tool providers is generally not realistic. This project instead uses:
- at-least-once delivery semantics;
- idempotency keys;
- persisted terminal results;
- explicit retry classification.

### Sync vs async

Low-latency tool calls remain synchronous. Long-running operations should move to an async execution path with durable queues and resumable state.

### Central policy vs local enforcement

A central policy model is easier to audit, but enforcement must happen in the execution path. This demo performs enforcement locally while keeping the policy interface replaceable.

## Repository layout

```text
cmd/gateway/              HTTP execution gateway
internal/executor/        request lifecycle and concurrency control
internal/policy/          capability policy engine
internal/trace/           audit/trace event helpers
ts-client/                TypeScript SDK/demo client
examples/                 example policy + request payloads
docs/                     architecture notes
.github/workflows/        CI
```

## Disclaimer

This is an independent proof-of-concept based only on public Metorial material and the public role description. It does not represent Metorial's internal architecture, source code, or infrastructure.
