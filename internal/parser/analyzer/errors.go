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

type ReservedKeywordError struct {
	key string
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

func NewReservedKeywordError(key string) error {
	return ReservedKeywordError{key}
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

func (e ReservedKeywordError) Error() string {
	return "reserved keyword " + e.key
}
