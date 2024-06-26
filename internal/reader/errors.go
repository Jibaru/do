package reader

type CanNotReadFileError struct {
	filename string
}

func NewCanNotReadFileError(filename string) error {
	return CanNotReadFileError{filename}
}

func (e CanNotReadFileError) Error() string {
	return "can not read file " + e.filename
}
