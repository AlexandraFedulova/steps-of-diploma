package graph

import (
	"math"
	"os"
	"testing"
)

func TestLouvainCommunities_twoCliques(t *testing.T) {
	g := New()
	// Кластер 1: A,B,C — плотные связи
	g.AddEdge("A", "B", 5)
	g.AddEdge("B", "A", 5)
	g.AddEdge("B", "C", 5)
	g.AddEdge("C", "A", 5)
	g.AddEdge("A", "C", 5)

	// Кластер 2: D,E,F
	g.AddEdge("D", "E", 5)
	g.AddEdge("E", "D", 5)
	g.AddEdge("E", "F", 5)
	g.AddEdge("F", "D", 5)
	g.AddEdge("D", "F", 5)

	res := g.LouvainCommunities()
	if len(res.Communities) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(res.Communities))
	}

	commOf := make(map[string]int)
	for _, a := range res.Assignments {
		commOf[a.NodeID] = a.CommunityID
	}

	if commOf["A"] != commOf["B"] || commOf["B"] != commOf["C"] {
		t.Errorf("cluster 1 split: A=%d B=%d C=%d", commOf["A"], commOf["B"], commOf["C"])
	}
	if commOf["D"] != commOf["E"] || commOf["E"] != commOf["F"] {
		t.Errorf("cluster 2 split: D=%d E=%d F=%d", commOf["D"], commOf["E"], commOf["F"])
	}
	if commOf["A"] == commOf["D"] {
		t.Error("expected different communities for cliques")
	}
}

func TestLouvainCommunities_everyNodeAssigned(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/health.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	res := g.LouvainCommunities()
	if len(res.Assignments) != g.NodeCount() {
		t.Fatalf("assignments = %d, nodes = %d", len(res.Assignments), g.NodeCount())
	}

	seen := make(map[string]bool)
	for _, a := range res.Assignments {
		if a.NodeID == "" {
			t.Fatal("empty node id")
		}
		if seen[a.NodeID] {
			t.Fatalf("duplicate assignment for %s", a.NodeID)
		}
		seen[a.NodeID] = true
	}

	covered := 0
	for _, c := range res.Communities {
		covered += len(c.Nodes)
	}
	if covered != g.NodeCount() {
		t.Fatalf("community nodes = %d, want %d", covered, g.NodeCount())
	}
}

func TestLouvainCommunities_negativeEdgesUseAbsWeight(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 5)
	g.AddEdge("B", "A", 5)
	g.AddEdge("C", "D", 5)
	g.AddEdge("D", "C", 5)
	g.AddEdge("B", "C", -10) // |вес| склеивает кластеры

	res := g.LouvainCommunities()
	if len(res.Communities) != 1 {
		t.Fatalf("expected 1 community (negative edge counts by |w|), got %d", len(res.Communities))
	}
}

func TestLouvainCommunities_singleNode(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})

	res := g.LouvainCommunities()
	if len(res.Assignments) != 1 || res.Assignments[0].CommunityID != 0 {
		t.Fatalf("got %+v", res.Assignments)
	}
}

func TestLouvainCommunities_emptyGraph(t *testing.T) {
	res := New().LouvainCommunities()
	if res == nil {
		t.Fatal("unexpected nil")
	}
	if len(res.Assignments) != 0 || len(res.Communities) != 0 {
		t.Fatalf("expected empty result, got assignments=%d communities=%d",
			len(res.Assignments), len(res.Communities))
	}
}

// Нет рёбер → m2 = 0 → тривиальный результат: у каждой вершины своё сообщество.
func TestLouvainCommunities_isolatedVerticesDistinctCommunities(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})
	g.AddNode(Node{ID: "B", Label: "B"})
	g.AddNode(Node{ID: "C", Label: "C"})

	res := g.LouvainCommunities()
	if len(res.Communities) != 3 {
		t.Fatalf("communities = %d, want 3", len(res.Communities))
	}
	seen := make(map[int]bool)
	for _, a := range res.Assignments {
		if seen[a.CommunityID] {
			t.Fatalf("duplicate community id %d for isolated nodes", a.CommunityID)
		}
		seen[a.CommunityID] = true
	}
	if math.IsNaN(res.Modularity) {
		t.Fatal("modularity is NaN")
	}
}

// На реальной матрице модульность конечна и алгоритм отрабатывает без сбоев.
func TestLouvainCommunities_sampleModularityIsFinite(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	res := g.LouvainCommunities()
	if math.IsNaN(res.Modularity) || math.IsInf(res.Modularity, 0) {
		t.Fatalf("Modularity = %v", res.Modularity)
	}
	if len(res.Assignments) != g.NodeCount() {
		t.Fatalf("assignments %d != nodes %d", len(res.Assignments), g.NodeCount())
	}
}
