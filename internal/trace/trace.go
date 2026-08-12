package trace

import (
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	TraceID  string         `json:"trace_id"`
	ActorID  string         `json:"actor_id"`
	AgentID  string         `json:"agent_id"`
	Tool     string         `json:"tool"`
	Op       string         `json:"operation"`
	Status   string         `json:"status"`
	At       time.Time      `json:"at"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func Emit(e Event) {
	e.At = time.Now().UTC()
	b, _ := json.Marshal(e)
	log.Printf("AUDIT %s", b)
}
