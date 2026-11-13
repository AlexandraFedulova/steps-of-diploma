package sociomatrix

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseCSV читает квадратную социоматрицу из CSV.
//
// Формат: первая строка — заголовки столбцов (первая ячейка пустая или совпадает с углом),
// далее строки «метка_строки;значения...». Число строк данных = числу столбцов данных.
//
// Значения ячеек:
//   - пусто, "0", "." — нет связи
//   - "+", "++", положительное число — положительная связь (вес = число или 1/2)
//   - "-", "--", отрицательное число — отрицательная связь
func ParseCSV(r io.Reader) (*Matrix, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}

	if len(records) < 2 {
		return nil, newParseError(0, 0, "matrix must have a header row and at least one data row")
	}

	header := records[0]
	if len(header) < 2 {
		return nil, newParseError(1, 0, "header must contain column labels")
	}

	colLabels := trimLabels(header[1:])
	if len(colLabels) == 0 {
		return nil, newParseError(1, 0, "no column labels")
	}
	if err := validateLabels(colLabels, 1, 2); err != nil {
		return nil, err
	}

	dataRows := records[1:]
	if len(dataRows) != len(colLabels) {
		return nil, newParseError(0, 0, fmt.Sprintf(
			"square matrix required: %d columns but %d data rows",
			len(colLabels), len(dataRows),
		))
	}

	rowLabels := make([]string, 0, len(dataRows))
	cells := make([][]Cell, 0, len(dataRows))

	for i, row := range dataRows {
		rowNum := i + 2
		expected := len(colLabels) + 1
		if len(row) != expected {
			return nil, newParseError(rowNum, 0, fmt.Sprintf(
				"expected %d fields, got %d",
				expected, len(row),
			))
		}

		rowLabel := strings.TrimSpace(row[0])
		if rowLabel == "" {
			return nil, newParseError(rowNum, 1, "empty row label")
		}
		rowLabels = append(rowLabels, rowLabel)

		rowCells := make([]Cell, len(colLabels))
		for j, raw := range row[1:] {
			cell, err := parseCell(raw)
			if err != nil {
				return nil, newParseError(rowNum, j+2, err.Error())
			}
			rowCells[j] = cell
		}
		cells = append(cells, rowCells)
	}

	if err := validateLabels(rowLabels, 0, 1); err != nil {
		return nil, err
	}

	if !sameLabels(rowLabels, colLabels) {
		return nil, newParseError(0, 0, "row labels must match column labels in the same order")
	}

	return &Matrix{Labels: rowLabels, Cells: cells}, nil
}

func trimLabels(labels []string) []string {
	out := make([]string, len(labels))
	for i, l := range labels {
		out[i] = strings.TrimSpace(l)
	}
	return out
}

func validateLabels(labels []string, row, colStart int) error {
	seen := make(map[string]struct{}, len(labels))
	for i, label := range labels {
		if label == "" {
			return newParseError(row, colStart+i, "empty label")
		}
		if _, dup := seen[label]; dup {
			return newParseError(row, colStart+i, fmt.Sprintf("duplicate label %q", label))
		}
		seen[label] = struct{}{}
	}
	return nil
}

func sameLabels(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func parseCell(raw string) (Cell, error) {
	s := strings.TrimSpace(raw)
	if s == "" || s == "0" || s == "." || s == "—" {
		return Cell{Kind: KindNone}, nil
	}

	switch s {
	case "+":
		return Cell{Kind: KindPositive, Weight: 1}, nil
	case "++":
		return Cell{Kind: KindPositive, Weight: 2}, nil
	case "-":
		return Cell{Kind: KindNegative, Weight: 1}, nil
	case "--":
		return Cell{Kind: KindNegative, Weight: 2}, nil
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Cell{}, fmt.Errorf("invalid cell value %q", raw)
	}
	if value == 0 {
		return Cell{Kind: KindNone}, nil
	}
	if value > 0 {
		return Cell{Kind: KindPositive, Weight: value}, nil
	}
	return Cell{Kind: KindNegative, Weight: -value}, nil
}
