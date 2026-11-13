package graph

// Component — компонента связности (слабая, для ориентированного графа).
type Component struct {
	ID    int      `json:"id"`
	Nodes []string `json:"nodes"`
}

// ConnectedComponents находит слабо связные компоненты графа.
//
// Рёбра рассматриваются как неориентированные: u → v соединяет u и v.
// Каждая вершина входит ровно в одну компоненту.
func (g *Graph) ConnectedComponents() []Component {
	if g.NodeCount() == 0 {
		return nil
	}

	undirected := g.undirectedAdjacency()
	visited := make(map[string]bool, g.NodeCount())
	components := make([]Component, 0)

	for _, start := range g.order {
		if visited[start] {
			continue
		}
		nodes := bfsUndirected(start, undirected, visited)
		components = append(components, Component{
			ID:    len(components),
			Nodes: nodes,
		})
	}

	return components
}

// ComponentCount возвращает число слабо связных компонент.
func (g *Graph) ComponentCount() int {
	return len(g.ConnectedComponents())
}
