package graph

import "github.com/mephi/sociometry/pkg/sociomatrix"

// FromMatrix строит ориентированный взвешенный граф из социоматрицы.
func FromMatrix(m *sociomatrix.Matrix) *Graph {
	g := New()

	for _, label := range m.Labels {
		g.AddNode(Node{ID: label, Label: label})
	}

	for i, row := range m.Cells {
		from := m.Labels[i]
		for j, cell := range row {
			if cell.Kind == sociomatrix.KindNone {
				continue
			}
			g.AddEdge(from, m.Labels[j], cell.SignedValue())
		}
	}

	return g
}
