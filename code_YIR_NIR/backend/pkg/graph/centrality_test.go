package graph

import (
	"os"
	"testing"
)

func TestInDegreeCentrality_weightedIncoming(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})
	g.AddNode(Node{ID: "B", Label: "B"})
	g.AddNode(Node{ID: "C", Label: "C"})

	g.AddEdge("A", "B", 2)  // входящее в B: +2
	g.AddEdge("C", "B", -1) // входящее в B: -1 → сумма 1
	g.AddEdge("B", "A", 3)  // входящее в A: +3

	scores := g.InDegreeCentralityMap()

	if scores["A"].InWeightSum != 3 {
		t.Errorf("A inWeight = %v, want 3", scores["A"].InWeightSum)
	}
	if scores["B"].InWeightSum != 1 {
		t.Errorf("B inWeight = %v, want 1", scores["B"].InWeightSum)
	}
	if scores["C"].InWeightSum != 0 {
		t.Errorf("C inWeight = %v, want 0", scores["C"].InWeightSum)
	}

	// n=3 → norm = 1/2
	if scores["A"].Value != 1.5 {
		t.Errorf("A value = %v, want 1.5", scores["A"].Value)
	}
	if scores["B"].Value != 0.5 {
		t.Errorf("B value = %v, want 0.5", scores["B"].Value)
	}
	if scores["C"].Value != 0 {
		t.Errorf("C value = %v, want 0", scores["C"].Value)
	}
}

func TestInDegreeCentrality_sampleCSV(t *testing.T) {
	f, err := os.Open("testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	scores := g.InDegreeCentrality()
	if len(scores) != 3 {
		t.Fatalf("scores count = %d, want 3", len(scores))
	}

	// sample: A→B:1, A→C:-1, B→A:1, B→C:-2, C→B:2
	want := map[string]float64{
		"A": 1,  // B→A
		"B": 3,  // A→B + C→B
		"C": -3, // A→C + B→C
	}

	for _, s := range scores {
		if s.InWeightSum != want[s.NodeID] {
			t.Errorf("node %s: inWeight = %v, want %v", s.NodeID, s.InWeightSum, want[s.NodeID])
		}
		if s.Value != want[s.NodeID]/2 {
			t.Errorf("node %s: value = %v, want %v", s.NodeID, s.Value, want[s.NodeID]/2)
		}
	}
}

func TestInDegreeCentrality_allNodesPresent(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/health.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	scores := g.InDegreeCentrality()
	if len(scores) != g.NodeCount() {
		t.Fatalf("scores = %d, nodes = %d", len(scores), g.NodeCount())
	}

	seen := make(map[string]bool)
	for _, s := range scores {
		if s.NodeID == "" {
			t.Fatal("empty node id in score")
		}
		seen[s.NodeID] = true
	}
	for _, n := range g.Nodes() {
		if !seen[n.ID] {
			t.Fatalf("missing centrality for node %s", n.ID)
		}
	}
}

// При n ≤ 1 нормировка не применяется: Value = 0 (см. InDegreeCentrality).
func TestInDegreeCentrality_singleNode_normalizedValueZero(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})

	scores := g.InDegreeCentrality()
	if len(scores) != 1 {
		t.Fatalf("len = %d", len(scores))
	}
	if scores[0].InWeightSum != 0 || scores[0].Value != 0 {
		t.Fatalf("got %+v, want InWeightSum=0 Value=0", scores[0])
	}
}

// Несколько рёбер в одну вершину суммируются (список смежности хранит каждое ребро отдельно).
func TestInDegreeCentrality_multipleEdgesToSameTarget(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})
	g.AddNode(Node{ID: "B", Label: "B"})
	g.AddEdge("A", "B", 1)
	g.AddEdge("A", "B", 2)

	s := g.InDegreeCentralityMap()["B"]
	if s.InWeightSum != 3 {
		t.Errorf("B InWeightSum = %v, want 3", s.InWeightSum)
	}
	// n=2 → norm=1
	if s.Value != 3 {
		t.Errorf("B Value = %v, want 3", s.Value)
	}
}

// Порядок scores совпадает с порядком добавления вершин (g.order).
func TestInDegreeCentrality_orderMatchesNodeOrder(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "Z", Label: "Z"})
	g.AddNode(Node{ID: "Y", Label: "Y"})
	g.AddEdge("Z", "Y", 1)

	scores := g.InDegreeCentrality()
	if len(scores) != 2 {
		t.Fatalf("len = %d", len(scores))
	}
	if scores[0].NodeID != "Z" || scores[1].NodeID != "Y" {
		t.Fatalf("order = %s, %s, want Z, Y", scores[0].NodeID, scores[1].NodeID)
	}
}
