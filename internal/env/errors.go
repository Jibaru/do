package env

import "fmt"

type CanNotReadFileError struct {
	Err string
}

func NewCanNotReadFileError(err string) *CanNotReadFileError {
	return &CanNotReadFileError{Err: err}
}

func (e *CanNotReadFileError) Error() string {
	return fmt.Sprintf("can not read file: %s", e.Err)
}
