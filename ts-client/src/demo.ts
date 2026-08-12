import { randomUUID } from "node:crypto";
import { ToolGatewayClient } from "./client.js";

const client = new ToolGatewayClient();

const result = await client.execute({
  actor_id: "user-123",
  agent_id: "support-agent",
  tool: "stripe",
  operation: "read",
  arguments: { customer_id: "cus_demo" },
  idempotency_key: randomUUID(),
  deadline_ms: 1500,
  trace_id: randomUUID(),
});

console.log(result);
