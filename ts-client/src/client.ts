export type ExecuteRequest = {
  actor_id: string;
  agent_id: string;
  tool: string;
  operation: string;
  arguments: Record<string, unknown>;
  idempotency_key: string;
  deadline_ms: number;
  trace_id: string;
};

export class ToolGatewayClient {
  constructor(private readonly baseUrl = "http://localhost:8080") {}

  async execute(req: ExecuteRequest) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), req.deadline_ms);
    try {
      const res = await fetch(`${this.baseUrl}/execute`, {
        method: "POST",
        headers: {"content-type": "application/json", "x-trace-id": req.trace_id},
        body: JSON.stringify(req),
        signal: controller.signal,
      });
      return await res.json();
    } finally {
      clearTimeout(timer);
    }
  }
}
