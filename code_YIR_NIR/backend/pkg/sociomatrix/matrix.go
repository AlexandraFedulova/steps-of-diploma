package sociomatrix

// Kind — тип связи в ячейке социоматрицы.
type Kind int

const (
	KindNone Kind = iota
	KindPositive
	KindNegative
)

func (k Kind) String() string {
	switch k {
	case KindPositive:
		return "positive"
	case KindNegative:
		return "negative"
	default:
		return "none"
	}
}

// Cell — значение одной ячейки (связь «от строки» → «к столбцу»).
type Cell struct {
	Kind   Kind
	Weight float64 // модуль веса; 0 при KindNone
}

// SignedValue возвращает знаковое значение: >0 положительная, <0 отрицательная, 0 — нет связи.
func (c Cell) SignedValue() float64 {
	switch c.Kind {
	case KindPositive:
		return c.Weight
	case KindNegative:
		return -c.Weight
	default:
		return 0
	}
}

// Edge — ненулевая связь в графе.
type Edge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Kind   string  `json:"kind"`
	Weight float64 `json:"weight"`
	Value  float64 `json:"value"`
}

// Matrix — квадратная социоматрица.
type Matrix struct {
	Labels []string `json:"labels"`
	Cells  [][]Cell `json:"cells"`
}

// Size возвращает размер матрицы (число участников).
func (m *Matrix) Size() int {
	return len(m.Labels)
}

// Edges возвращает все ненулевые связи (без диагонали по умолчанию не фильтруем).
func (m *Matrix) Edges() []Edge {
	var edges []Edge
	for i, row := range m.Cells {
		for j, cell := range row {
			if cell.Kind == KindNone {
				continue
			}
			edges = append(edges, Edge{
				From:   m.Labels[i],
				To:     m.Labels[j],
				Kind:   cell.Kind.String(),
				Weight: cell.Weight,
				Value:  cell.SignedValue(),
			})
		}
	}
	return edges
}
