package graph

import (
	"os"
	"testing"
)

func TestConnectedComponents_healthCSV(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/health.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	components := g.ConnectedComponents()
	if len(components) == 0 {
		t.Fatal("expected at least one component")
	}

	total := 0
	for _, c := range components {
		total += len(c.Nodes)
	}
	if total != g.NodeCount() {
		t.Fatalf("nodes in components = %d, graph nodes = %d", total, g.NodeCount())
	}
}
