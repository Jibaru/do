package reader

type CanNotReadFileError struct {
	err error
}

func NewCanNotReadFileError(err error) error {
	return CanNotReadFileError{err}
}

func (e CanNotReadFileError) Error() string {
	return "can not read file: " + e.err.Error()
}
