package analyzer

import "fmt"

type ReadingExpressionError struct {
	parts []string
}

type RepeatedKeyError struct {
	key string
}

type InvalidValueError struct {
	value string
}

func NewReadingExpressionError(parts ...string) error {
	return ReadingExpressionError{parts}
}

func NewRepeatedKeyError(key string) error {
	return RepeatedKeyError{key}
}

func NewInvalidValueError(value string) error {
	return InvalidValueError{value}
}

func (e ReadingExpressionError) Error() string {
	return "error reading expression: " + fmt.Sprintf("%v", e.parts)
}

func (e RepeatedKeyError) Error() string {
	return "repeated key " + e.key
}

func (e InvalidValueError) Error() string {
	return "invalid value " + e.value
}
