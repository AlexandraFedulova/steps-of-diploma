package graph

import "fmt"

// BFSResult — результат обхода в ширину.
type BFSResult struct {
	Start    string         `json:"start"`
	Order    []string       `json:"order"`
	Distance map[string]int `json:"distance"`
}

// BFS выполняет обход в ширину от start по направлению рёбер (исходящие).
func (g *Graph) BFS(start string) (*BFSResult, error) {
	if _, ok := g.nodes[start]; !ok {
		return nil, fmt.Errorf("unknown node %q", start)
	}

	visited := make(map[string]bool, g.NodeCount())
	distance := make(map[string]int, g.NodeCount())
	order := make([]string, 0, g.NodeCount())

	queue := []string{start}
	visited[start] = true
	distance[start] = 0

	for head := 0; head < len(queue); head++ {
		v := queue[head]
		order = append(order, v)

		for _, e := range g.adj[v] {
			if visited[e.To] {
				continue
			}
			visited[e.To] = true
			distance[e.To] = distance[v] + 1
			queue = append(queue, e.To)
		}
	}

	return &BFSResult{
		Start:    start,
		Order:    order,
		Distance: distance,
	}, nil
}

// bfsUndirected обходит неориентированный граф и возвращает вершины компоненты.
func bfsUndirected(start string, adj map[string][]string, visited map[string]bool) []string {
	queue := []string{start}
	visited[start] = true
	component := []string{start}

	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, next := range adj[v] {
			if visited[next] {
				continue
			}
			visited[next] = true
			queue = append(queue, next)
			component = append(component, next)
		}
	}

	return component
}

// undirectedAdjacency строит список смежности без учёта направления рёбер.
func (g *Graph) undirectedAdjacency() map[string][]string {
	adj := make(map[string][]string, len(g.nodes))
	for id := range g.nodes {
		adj[id] = nil
	}

	for from, edges := range g.adj {
		for _, e := range edges {
			adj[from] = append(adj[from], e.To)
			adj[e.To] = append(adj[e.To], from)
		}
	}

	return adj
}
