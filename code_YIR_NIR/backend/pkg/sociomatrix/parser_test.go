package sociomatrix

import (
	"errors"
	"strings"
	"testing"
)

// Корректная квадратная матрица: метки, размер, положительные и отрицательные ячейки, рёбра.
func TestParseCSV_correct(t *testing.T) {
	csv := `,A,B,C
A,0,+,-1
B,2,0,--
C,.,-3,0
`

	m, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := []string{"A", "B", "C"}; len(m.Labels) != 3 || m.Labels[0] != got[0] || m.Labels[1] != got[1] || m.Labels[2] != got[2] {
		t.Fatalf("labels = %v, want %v", m.Labels, got)
	}

	if m.Size() != 3 {
		t.Fatalf("size = %d, want 3", m.Size())
	}

	// A→B: +
	if m.Cells[0][1].Kind != KindPositive || m.Cells[0][1].Weight != 1 || m.Cells[0][1].SignedValue() != 1 {
		t.Errorf("A→B: %+v", m.Cells[0][1])
	}
	// A→C: -1 (numeric negative)
	if m.Cells[0][2].Kind != KindNegative || m.Cells[0][2].Weight != 1 || m.Cells[0][2].SignedValue() != -1 {
		t.Errorf("A→C: %+v", m.Cells[0][2])
	}
	// B→A: 2
	if m.Cells[1][0].Kind != KindPositive || m.Cells[1][0].SignedValue() != 2 {
		t.Errorf("B→A: %+v", m.Cells[1][0])
	}
	// B→C: --
	if m.Cells[1][2].Kind != KindNegative || m.Cells[1][2].Weight != 2 || m.Cells[1][2].SignedValue() != -2 {
		t.Errorf("B→C: %+v", m.Cells[1][2])
	}
	// C→B: -3
	if m.Cells[2][1].Kind != KindNegative || m.Cells[2][1].Weight != 3 || m.Cells[2][1].SignedValue() != -3 {
		t.Errorf("C→B: %+v", m.Cells[2][1])
	}

	edges := m.Edges()
	if len(edges) != 5 {
		t.Fatalf("edges = %d, want 5: %+v", len(edges), edges)
	}
}

func TestParseCSV_valid(t *testing.T) {
	csv := `,A,B,C
A,0,+,-
B,1,0,-2
C,.,++,0
`

	m, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Size() != 3 {
		t.Fatalf("size = %d, want 3", m.Size())
	}

	if m.Cells[0][1].Kind != KindPositive || m.Cells[0][1].Weight != 1 {
		t.Errorf("A→B: got %+v", m.Cells[0][1])
	}
	if m.Cells[0][2].Kind != KindNegative || m.Cells[0][2].Weight != 1 {
		t.Errorf("A→C: got %+v", m.Cells[0][2])
	}
	if m.Cells[1][2].Kind != KindNegative || m.Cells[1][2].Weight != 2 {
		t.Errorf("B→C: got %+v", m.Cells[1][2])
	}
	if m.Cells[2][1].Kind != KindPositive || m.Cells[2][1].Weight != 2 {
		t.Errorf("C→B: got %+v", m.Cells[2][1])
	}

	if len(m.Edges()) != 5 {
		t.Fatalf("edges count = %d, want 5", len(m.Edges()))
	}
}

// Некорректный синтаксис CSV (ошибка encoding/csv до валидации матрицы).
func TestParseCSV_incorrectSyntax(t *testing.T) {
	// Незакрытая кавычка во второй строке
	bad := ",A,B\n\"broken"
	_, err := ParseCSV(strings.NewReader(bad))
	if err == nil {
		t.Fatal("expected error for broken CSV quoting")
	}
	if !strings.Contains(err.Error(), "read csv") {
		t.Fatalf("expected wrapped read error, got: %v", err)
	}
}

// Пустой поток или недостаточно строк (нет строки данных).
func TestParseCSV_emptyFile(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{"empty", ""},
		{"only_spaces", "   \n  \t  "},
		{"header_only", ",A,B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCSV(strings.NewReader(tt.raw))
			if err == nil {
				t.Fatal("expected error")
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("expected *ParseError, got %T: %v", err, err)
			}
		})
	}
}

func TestParseCSV_negativeWeights(t *testing.T) {
	csv := `,P,Q,R
P,0,-,-5.5
Q,-2,0,+
R,--,1,0
`
	m, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		from, to int // индексы меток P=0,Q=1,R=2
		want     float64
	}{
		{0, 1, -1},   // P→Q: "-"
		{0, 2, -5.5}, // P→R
		{1, 0, -2},   // Q→P
		{2, 0, -2},   // R→P: "--"
		{1, 2, 1},    // Q→R: "+"
	}
	for _, c := range cases {
		cell := m.Cells[c.from][c.to]
		if got := cell.SignedValue(); got != c.want {
			t.Errorf("%s→%s: signed=%v want %v (cell=%+v)",
				m.Labels[c.from], m.Labels[c.to], got, c.want, cell)
		}
		if cell.SignedValue() < 0 && cell.Kind != KindNegative {
			t.Errorf("%s→%s: expected KindNegative for negative signed value", m.Labels[c.from], m.Labels[c.to])
		}
	}
}

func TestParseCSV_notSquare(t *testing.T) {
	csv := `,A,B
A,0,1
`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected error for non-square matrix")
	}
}

func TestParseCSV_mismatchedRowWidth(t *testing.T) {
	csv := `,A,B
A,0
B,1,0
`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected error for row width mismatch")
	}
}

func TestParseCSV_labelMismatch(t *testing.T) {
	csv := `,A,B
X,0,1
Y,1,0
`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected error for label mismatch")
	}
}

func TestParseCSV_duplicateLabels(t *testing.T) {
	csv := `,A,A
A,0,1
A,1,0
`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected error for duplicate column labels")
	}
}

func TestParseCSV_invalidCell(t *testing.T) {
	csv := `,A,B
A,0,foo
B,1,0
`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected parse error")
	}

	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Row != 2 || pe.Column != 3 {
		t.Errorf("position = row %d col %d", pe.Row, pe.Column)
	}
}

func TestParseCell_numeric(t *testing.T) {
	tests := []struct {
		in     string
		kind   Kind
		weight float64
		signed float64
	}{
		{"2.5", KindPositive, 2.5, 2.5},
		{"-3", KindNegative, 3, -3},
		{"0", KindNone, 0, 0},
		{"+", KindPositive, 1, 1},
		{"--", KindNegative, 2, -2},
	}

	for _, tt := range tests {
		cell, err := parseCell(tt.in)
		if err != nil {
			t.Fatalf("%q: %v", tt.in, err)
		}
		if cell.Kind != tt.kind || cell.Weight != tt.weight || cell.SignedValue() != tt.signed {
			t.Errorf("%q: got kind=%v weight=%v signed=%v", tt.in, cell.Kind, cell.Weight, cell.SignedValue())
		}
	}
}
