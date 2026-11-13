package sociomatrix

import "fmt"

// ParseError описывает ошибку разбора или валидации CSV.
type ParseError struct {
	Row     int    // 1-based, 0 если не привязана к строке
	Column  int    // 1-based, 0 если не привязана к столбцу
	Message string
}

func (e *ParseError) Error() string {
	switch {
	case e.Row > 0 && e.Column > 0:
		return fmt.Sprintf("row %d, column %d: %s", e.Row, e.Column, e.Message)
	case e.Row > 0:
		return fmt.Sprintf("row %d: %s", e.Row, e.Message)
	default:
		return e.Message
	}
}

func newParseError(row, col int, msg string) error {
	return &ParseError{Row: row, Column: col, Message: msg}
}
