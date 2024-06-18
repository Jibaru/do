package analyzer

import "fmt"

type ReadingExpressionError struct {
	parts []string
}

func NewReadingExpressionError(parts ...string) error {
	return ReadingExpressionError{parts}
}

func (e ReadingExpressionError) Error() string {
	return "error reading expression: " + fmt.Sprintf("%v", e.parts)
}
