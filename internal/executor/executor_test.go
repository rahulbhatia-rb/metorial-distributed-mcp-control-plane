package executor

import (
	"context"
	"testing"

	"github.com/rahulbhatia-rb/metorial-distributed-mcp-control-plane/internal/policy"
)

func TestIdempotentReplay(t *testing.T) {
	p := policy.NewStatic(map[string][]string{"a": {"stripe.read"}})
	e := New(p, 1)
	req := Request{ActorID: "u", AgentID: "a", Tool: "stripe", Operation: "read", IdempotencyKey: "k"}
	a := e.Execute(context.Background(), req)
	b := e.Execute(context.Background(), req)
	if a.Status != "completed" || b.Status != "completed" {
		t.Fatalf("expected completed replays")
	}
}

func TestPolicyDeny(t *testing.T) {
	p := policy.NewStatic(map[string][]string{"a": {"stripe.read"}})
	e := New(p, 1)
	r := e.Execute(context.Background(), Request{
		ActorID: "u", AgentID: "a", Tool: "stripe", Operation: "write", IdempotencyKey: "k2",
	})
	if r.Status != "denied" {
		t.Fatalf("expected denied, got %s", r.Status)
	}
}
