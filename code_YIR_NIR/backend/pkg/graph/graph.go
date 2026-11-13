package graph

// Node — вершина ориентированного взвешенного графа.
type Node struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Edge — исходящее ребро в списке смежности (from → To с весом Weight).
type Edge struct {
	To     string  `json:"to"`
	Weight float64 `json:"weight"`
}

// Graph — ориентированный взвешенный граф (adjacency list).
type Graph struct {
	Directed bool `json:"directed"`

	nodes map[string]Node
	adj   map[string][]Edge
	order []string
}

// New создаёт пустой ориентированный граф.
func New() *Graph {
	return &Graph{
		Directed: true,
		nodes:    make(map[string]Node),
		adj:      make(map[string][]Edge),
	}
}

// AddNode добавляет вершину. Повторный вызов с тем же ID обновляет метку.
func (g *Graph) AddNode(n Node) {
	if n.Label == "" {
		n.Label = n.ID
	}
	if _, exists := g.nodes[n.ID]; !exists {
		g.order = append(g.order, n.ID)
	}
	g.nodes[n.ID] = n
	if _, ok := g.adj[n.ID]; !ok {
		g.adj[n.ID] = nil
	}
}

// AddEdge добавляет ориентированное ребро from → to с весом.
func (g *Graph) AddEdge(from, to string, weight float64) {
	if _, ok := g.nodes[from]; !ok {
		g.AddNode(Node{ID: from, Label: from})
	}
	if _, ok := g.nodes[to]; !ok {
		g.AddNode(Node{ID: to, Label: to})
	}
	g.adj[from] = append(g.adj[from], Edge{To: to, Weight: weight})
}

// NodeCount возвращает число вершин.
func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

// EdgeCount возвращает число рёбер.
func (g *Graph) EdgeCount() int {
	n := 0
	for _, edges := range g.adj {
		n += len(edges)
	}
	return n
}

// Nodes возвращает вершины в порядке добавления.
func (g *Graph) Nodes() []Node {
	out := make([]Node, 0, len(g.order))
	for _, id := range g.order {
		out = append(out, g.nodes[id])
	}
	return out
}

// Neighbors возвращает исходящие рёбра вершины id.
func (g *Graph) Neighbors(id string) []Edge {
	edges := g.adj[id]
	if len(edges) == 0 {
		return nil
	}
	out := make([]Edge, len(edges))
	copy(out, edges)
	return out
}

// AdjacencyList возвращает копию списка смежности.
func (g *Graph) AdjacencyList() map[string][]Edge {
	out := make(map[string][]Edge, len(g.adj))
	for from, edges := range g.adj {
		if len(edges) == 0 {
			out[from] = nil
			continue
		}
		cp := make([]Edge, len(edges))
		copy(cp, edges)
		out[from] = cp
	}
	return out
}
