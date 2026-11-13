package graph

// EdgeDTO — ребро для JSON payload (source → target).
type EdgeDTO struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
}

// WeightMatrixDTO — веса социоматрицы (строка i → столбец j).
type WeightMatrixDTO struct {
	Labels []string    `json:"labels"`
	Values [][]float64 `json:"values"`
}

// GraphPayload — полный граф для frontend.
type GraphPayload struct {
	Directed   bool              `json:"directed"`
	NodeCount  int               `json:"nodeCount"`
	EdgeCount  int               `json:"edgeCount"`
	Nodes      []Node            `json:"nodes"`
	Edges      []EdgeDTO         `json:"edges"`
	Centrality []CentralityScore `json:"centrality"`
	Community  []NodeCommunity   `json:"community"`
	Weights    WeightMatrixDTO   `json:"weights"`
}

// Payload собирает DTO с узлами, рёбрами, centrality, кликовыми сообществами и матрицей весов.
func (g *Graph) Payload() *GraphPayload {
	return &GraphPayload{
		Directed:   g.Directed,
		NodeCount:  g.NodeCount(),
		EdgeCount:  g.EdgeCount(),
		Nodes:      g.Nodes(),
		Edges:      g.edgesDTO(),
		Centrality: g.InDegreeCentrality(),
		Community:  g.CliqueCommunityAssignments(),
		Weights:    g.weightMatrixDTO(),
	}
}

func (g *Graph) edgesDTO() []EdgeDTO {
	out := make([]EdgeDTO, 0, g.EdgeCount())
	for _, from := range g.order {
		for _, e := range g.adj[from] {
			out = append(out, EdgeDTO{
				Source: from,
				Target: e.To,
				Weight: e.Weight,
			})
		}
	}
	return out
}

func (g *Graph) weightMatrixDTO() WeightMatrixDTO {
	n := len(g.order)
	labels := make([]string, n)
	copy(labels, g.order)

	index := make(map[string]int, n)
	for i, id := range g.order {
		index[id] = i
	}

	values := make([][]float64, n)
	for i := range values {
		values[i] = make([]float64, n)
	}

	for from, edges := range g.adj {
		i := index[from]
		for _, e := range edges {
			values[i][index[e.To]] = e.Weight
		}
	}

	return WeightMatrixDTO{Labels: labels, Values: values}
}
