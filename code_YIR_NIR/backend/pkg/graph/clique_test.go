package graph

import (
	"os"
	"testing"
)

func TestCliqueCommunities_frgm(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/frgm.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	res := g.CliqueCommunities()
	if len(res.MutualCliques) != 4 {
		t.Fatalf("mutual cliques = %d, want 4", len(res.MutualCliques))
	}
	if len(res.Strict) != 0 {
		t.Fatalf("expected no strict closed cliques in frgm, got %d", len(res.Strict))
	}
}

func TestCliqueCommunities_userExample(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 2)
	g.AddEdge("A", "C", 1)
	g.AddEdge("B", "A", 1)
	g.AddEdge("B", "C", 2)
	g.AddEdge("C", "B", 2)

	res := g.CliqueCommunities()
	found := false
	for _, c := range res.AnyPositiveCliques {
		if c.Size == 3 {
			found = true
		}
	}
	if !found {
		t.Fatal("expected any-positive clique of size 3 (A,B,C)")
	}
}
