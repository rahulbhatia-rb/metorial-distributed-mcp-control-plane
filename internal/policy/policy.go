package policy

import "fmt"

type Engine interface {
	Allow(agentID, capability string) error
}

type Static struct {
	grants map[string]map[string]struct{}
}

func NewStatic(in map[string][]string) *Static {
	g := map[string]map[string]struct{}{}
	for agent, caps := range in {
		g[agent] = map[string]struct{}{}
		for _, c := range caps {
			g[agent][c] = struct{}{}
		}
	}
	return &Static{grants: g}
}

func (s *Static) Allow(agentID, capability string) error {
	if _, ok := s.grants[agentID][capability]; !ok {
		return fmt.Errorf("policy denied %s for %s", capability, agentID)
	}
	return nil
}
