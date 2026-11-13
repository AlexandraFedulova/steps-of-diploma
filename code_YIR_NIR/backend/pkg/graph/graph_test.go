package graph

import (
	"os"
	"testing"
)

func TestParseCSV_sample(t *testing.T) {
	f, err := os.Open("testdata/sample.csv")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatalf("ParseCSV: %v", err)
	}

	if !g.Directed {
		t.Error("graph must be directed")
	}
	if g.NodeCount() != 3 {
		t.Fatalf("nodes = %d, want 3", g.NodeCount())
	}
	if g.EdgeCount() != 5 {
		t.Fatalf("edges = %d, want 5", g.EdgeCount())
	}

	adj := g.AdjacencyList()
	if len(adj["A"]) != 2 {
		t.Fatalf("A out-degree = %d, want 2", len(adj["A"]))
	}
}

func TestParseCSV_health(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/health.csv")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatalf("ParseCSV: %v", err)
	}

	if g.NodeCount() != 15 {
		t.Fatalf("nodes = %d, want 15", g.NodeCount())
	}
	if g.EdgeCount() == 0 {
		t.Fatal("expected non-empty graph")
	}

	for _, n := range g.Nodes() {
		if n.ID == "" || n.Label == "" {
			t.Fatalf("invalid node: %+v", n)
		}
		_ = g.Neighbors(n.ID)
	}
}

func TestAdjacencyList_directedWeighted(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})
	g.AddNode(Node{ID: "B", Label: "B"})
	g.AddEdge("A", "B", 1.5)
	g.AddEdge("B", "A", -2)

	ab := g.Neighbors("A")
	if len(ab) != 1 || ab[0].To != "B" || ab[0].Weight != 1.5 {
		t.Fatalf("A→B: %+v", ab)
	}

	ba := g.Neighbors("B")
	if len(ba) != 1 || ba[0].Weight != -2 {
		t.Fatalf("B→A: %+v", ba)
	}
}
