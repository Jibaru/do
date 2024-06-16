package analyzer

import "fmt"

type ReadingExpressionError struct {
	parts []string
}

type CanNotParseMapValueError struct {
	err error
}

type CanNotParseBoolValueError struct {
	err error
}

func NewReadingExpressionError(parts ...string) error {
	return ReadingExpressionError{parts}
}

func NewCanNotParseMapValueError(err error) error {
	return CanNotParseMapValueError{err}
}

func NewCanNotParseBoolValueError(err error) error {
	return CanNotParseBoolValueError{err}
}

func (e ReadingExpressionError) Error() string {
	return "error reading expression: " + fmt.Sprintf("%v", e.parts)
}

func (e CanNotParseMapValueError) Error() string {
	return "can not parse map value: " + e.err.Error()
}

func (e CanNotParseBoolValueError) Error() string {
	return "can not parse bool value: " + e.err.Error()
}
