package graph

import (
	"strings"
	"testing"
)

func TestGraphPayload(t *testing.T) {
	csv := `,A,B,C
A,0,1,-1
B,0,0,1
C,1,0,0`

	g, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	p := g.Payload()

	if len(p.Nodes) != 3 {
		t.Fatalf("nodes: got %d", len(p.Nodes))
	}
	if len(p.Edges) != 4 {
		t.Fatalf("edges: got %d", len(p.Edges))
	}
	if len(p.Centrality) != 3 {
		t.Fatalf("centrality: got %d", len(p.Centrality))
	}
	if len(p.Community) != 3 {
		t.Fatalf("community: got %d", len(p.Community))
	}
	if len(p.Weights.Labels) != 3 || len(p.Weights.Values) != 3 {
		t.Fatalf("weights matrix: labels=%d values=%d", len(p.Weights.Labels), len(p.Weights.Values))
	}

	if p.Weights.Values[0][1] != 1 {
		t.Errorf("A→B weight: got %v", p.Weights.Values[0][1])
	}
	if p.Edges[0].Source == "" || p.Edges[0].Target == "" {
		t.Errorf("edge dto missing source/target: %+v", p.Edges[0])
	}

	commNodes := make(map[string]bool)
	for _, a := range p.Community {
		commNodes[a.NodeID] = true
	}
	for _, n := range p.Nodes {
		if !commNodes[n.ID] {
			t.Errorf("no community for node %s", n.ID)
		}
	}
}
