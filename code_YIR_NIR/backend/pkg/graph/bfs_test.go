package graph

import (
	"reflect"
	"testing"
)

func TestBFS_directed(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 1)
	g.AddEdge("B", "C", 1)
	g.AddNode(Node{ID: "D", Label: "D"}) // изолирована от A по направлению рёбер

	res, err := g.BFS("A")
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Order) != 3 {
		t.Fatalf("reachable = %d, want 3: %v", len(res.Order), res.Order)
	}
	if res.Distance["C"] != 2 {
		t.Errorf("dist C = %d, want 2", res.Distance["C"])
	}
	if _, ok := res.Distance["D"]; ok {
		t.Error("D should not be reachable from A")
	}

	resC, err := g.BFS("C")
	if err != nil {
		t.Fatal(err)
	}
	if len(resC.Order) != 1 {
		t.Fatalf("from C reachable = %d, want 1", len(resC.Order))
	}
}

func TestBFS_unknownStart(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})

	if _, err := g.BFS("X"); err == nil {
		t.Fatal("expected error for unknown start node")
	}
}

func TestBFS_singleNode(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})

	res, err := g.BFS("A")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res.Order, []string{"A"}) {
		t.Fatalf("order = %v", res.Order)
	}
	if res.Distance["A"] != 0 {
		t.Fatalf("distance A = %d", res.Distance["A"])
	}
}

// Ромб: A на уровне 0, B и C на 1, D на 2. Веса рёбер для BFS не важны.
func TestBFS_diamondGraph_distances(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 1)
	g.AddEdge("A", "C", 1)
	g.AddEdge("B", "D", 1)
	g.AddEdge("C", "D", 1)

	res, err := g.BFS("A")
	if err != nil {
		t.Fatal(err)
	}
	wantDist := map[string]int{"A": 0, "B": 1, "C": 1, "D": 2}
	if !reflect.DeepEqual(res.Distance, wantDist) {
		t.Errorf("distance = %#v, want %#v", res.Distance, wantDist)
	}
	if len(res.Order) != 4 {
		t.Fatalf("reachable count = %d, order=%v", len(res.Order), res.Order)
	}
	if res.Order[0] != "A" || res.Distance[res.Order[len(res.Order)-1]] != 2 {
		t.Fatalf("unexpected order %v", res.Order)
	}
}

// Направленный цикл: из каждой вершины достижимы все три.
func TestBFS_directedCycle_reachesAll(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 1)
	g.AddEdge("B", "C", 1)
	g.AddEdge("C", "A", 1)

	res, err := g.BFS("A")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Order) != 3 {
		t.Fatalf("order len = %d, order=%v", len(res.Order), res.Order)
	}
	if res.Distance["A"] != 0 || res.Distance["B"] != 1 || res.Distance["C"] != 2 {
		t.Fatalf("distances = %#v", res.Distance)
	}
}

func TestConnectedComponents_twoComponents(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 1)
	g.AddEdge("C", "D", 1)

	components := g.ConnectedComponents()
	if len(components) != 2 {
		t.Fatalf("components = %d, want 2", len(components))
	}

	sizes := map[int]int{}
	for _, c := range components {
		sizes[len(c.Nodes)]++
	}
	if sizes[2] != 2 {
		t.Fatalf("expected two components of size 2, got sizes %v", sizes)
	}
}

func TestConnectedComponents_singleDirectedEdge(t *testing.T) {
	g := New()
	g.AddEdge("A", "B", 1)

	components := g.ConnectedComponents()
	if len(components) != 1 || len(components[0].Nodes) != 2 {
		t.Fatalf("got %+v", components)
	}
}

func TestConnectedComponents_isolatedNodes(t *testing.T) {
	g := New()
	g.AddNode(Node{ID: "A", Label: "A"})
	g.AddNode(Node{ID: "B", Label: "B"})
	g.AddNode(Node{ID: "C", Label: "C"})

	components := g.ConnectedComponents()
	if len(components) != 3 {
		t.Fatalf("components = %d, want 3", len(components))
	}
}

func TestConnectedComponents_coversAllNodes(t *testing.T) {
	g := buildSampleGraph(t)

	components := g.ConnectedComponents()
	seen := make(map[string]int)
	for _, c := range components {
		for _, id := range c.Nodes {
			seen[id]++
		}
	}

	if len(seen) != g.NodeCount() {
		t.Fatalf("covered %d nodes, want %d", len(seen), g.NodeCount())
	}
	for id, count := range seen {
		if count != 1 {
			t.Fatalf("node %s in %d components", id, count)
		}
	}
}

func buildSampleGraph(t *testing.T) *Graph {
	t.Helper()
	g := New()
	g.AddEdge("A", "B", 1)
	g.AddEdge("B", "C", 1)
	g.AddEdge("X", "Y", 1)
	return g
}
