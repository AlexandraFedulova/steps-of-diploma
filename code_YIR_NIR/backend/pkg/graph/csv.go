package graph

import (
	"io"

	"github.com/mephi/sociometry/pkg/sociomatrix"
)

// ParseCSV читает CSV социоматрицы и возвращает граф.
func ParseCSV(r io.Reader) (*Graph, error) {
	matrix, err := sociomatrix.ParseCSV(r)
	if err != nil {
		return nil, err
	}
	return FromMatrix(matrix), nil
}
