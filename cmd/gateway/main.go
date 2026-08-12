package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rahulbhatia-rb/metorial-distributed-mcp-control-plane/internal/executor"
	"github.com/rahulbhatia-rb/metorial-distributed-mcp-control-plane/internal/policy"
)

func main() {
	p := policy.NewStatic(map[string][]string{
		"support-agent": {"stripe.read", "sentry.read", "zendesk.write"},
		"ops-agent":     {"aws.read", "consul.read", "vault.read"},
	})
	ex := executor.New(p, 32)

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		var req executor.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		if req.DeadlineMS > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(req.DeadlineMS)*time.Millisecond)
			defer cancel()
		}
		res := ex.Execute(ctx, req)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
