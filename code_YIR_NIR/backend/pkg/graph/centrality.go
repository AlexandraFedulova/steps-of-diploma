package graph

// CentralityScore — degree centrality вершины по входящим взвешенным связям.
type CentralityScore struct {
	NodeID      string  `json:"nodeId"`
	InWeightSum float64 `json:"inWeightSum"` // сумма весов входящих рёбер
	Value       float64 `json:"value"`       // нормализованная centrality: inWeightSum / (n-1)
}

// InDegreeCentrality вычисляет degree centrality для каждой вершины.
//
// Учитываются только входящие рёбра (u → v добавляет вес к v).
// Веса суммируются со знаком (отрицательные связи уменьшают centrality).
// Value нормализуется делением на (n-1); при n ≤ 1 Value = 0.
func (g *Graph) InDegreeCentrality() []CentralityScore {
	n := g.NodeCount()
	inWeight := make(map[string]float64, n)

	for _, id := range g.order {
		inWeight[id] = 0
	}

	for _, edges := range g.adj {
		for _, e := range edges {
			inWeight[e.To] += e.Weight
		}
	}

	norm := 0.0
	if n > 1 {
		norm = 1.0 / float64(n-1)
	}

	scores := make([]CentralityScore, 0, len(g.order))
	for _, id := range g.order {
		sum := inWeight[id]
		scores = append(scores, CentralityScore{
			NodeID:      id,
			InWeightSum: sum,
			Value:       sum * norm,
		})
	}

	return scores
}

// InDegreeCentralityMap возвращает centrality как map[nodeID]CentralityScore.
func (g *Graph) InDegreeCentralityMap() map[string]CentralityScore {
	scores := g.InDegreeCentrality()
	out := make(map[string]CentralityScore, len(scores))
	for _, s := range scores {
		out[s.NodeID] = s
	}
	return out
}
