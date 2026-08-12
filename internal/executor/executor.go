package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/rahulbhatia-rb/metorial-distributed-mcp-control-plane/internal/policy"
	"github.com/rahulbhatia-rb/metorial-distributed-mcp-control-plane/internal/trace"
)

type Request struct {
	ActorID        string         `json:"actor_id"`
	AgentID        string         `json:"agent_id"`
	Tool           string         `json:"tool"`
	Operation      string         `json:"operation"`
	Arguments      map[string]any `json:"arguments"`
	IdempotencyKey string         `json:"idempotency_key"`
	DeadlineMS     int            `json:"deadline_ms"`
	TraceID        string         `json:"trace_id"`
}

type Result struct {
	Status    string         `json:"status"`
	Retryable bool           `json:"retryable"`
	Output    map[string]any `json:"output,omitempty"`
	Error     string         `json:"error,omitempty"`
}

type Executor struct {
	policy policy.Engine
	sem    chan struct{}
	mu     sync.Mutex
	done   map[string]Result
}

func New(p policy.Engine, maxConcurrent int) *Executor {
	return &Executor{
		policy: p,
		sem:    make(chan struct{}, maxConcurrent),
		done:   map[string]Result{},
	}
}

func (e *Executor) Execute(ctx context.Context, req Request) Result {
	if req.ActorID == "" || req.AgentID == "" || req.Tool == "" || req.Operation == "" {
		return Result{Status: "rejected", Error: "missing identity or tool fields"}
	}
	if req.IdempotencyKey == "" {
		return Result{Status: "rejected", Error: "idempotency_key required"}
	}

	e.mu.Lock()
	if cached, ok := e.done[req.IdempotencyKey]; ok {
		e.mu.Unlock()
		return cached
	}
	e.mu.Unlock()

	capability := fmt.Sprintf("%s.%s", req.Tool, req.Operation)
	if err := e.policy.Allow(req.AgentID, capability); err != nil {
		trace.Emit(trace.Event{TraceID: req.TraceID, ActorID: req.ActorID, AgentID: req.AgentID, Tool: req.Tool, Op: req.Operation, Status: "denied"})
		return Result{Status: "denied", Error: err.Error()}
	}

	select {
	case e.sem <- struct{}{}:
		defer func() { <-e.sem }()
	case <-ctx.Done():
		return Result{Status: "timeout", Retryable: true, Error: ctx.Err().Error()}
	}

	trace.Emit(trace.Event{TraceID: req.TraceID, ActorID: req.ActorID, AgentID: req.AgentID, Tool: req.Tool, Op: req.Operation, Status: "started"})

	res := Result{
		Status: "completed",
		Output: map[string]any{
			"tool":      req.Tool,
			"operation": req.Operation,
			"echo":      req.Arguments,
		},
	}

	e.mu.Lock()
	e.done[req.IdempotencyKey] = res
	e.mu.Unlock()

	trace.Emit(trace.Event{TraceID: req.TraceID, ActorID: req.ActorID, AgentID: req.AgentID, Tool: req.Tool, Op: req.Operation, Status: "completed"})
	return res
}
