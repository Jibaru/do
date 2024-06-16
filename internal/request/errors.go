package request

type CanNotDoRequestError struct {
	err error
}

type CanNotReadResponseBodyError struct {
	err error
}

func NewCanNotDoRequestError(err error) error {
	return CanNotDoRequestError{err}
}

func NewCanNotReadResponseBodyError(err error) error {
	return CanNotReadResponseBodyError{err}
}

func (e CanNotDoRequestError) Error() string {
	return "can not do request: " + e.err.Error()
}

func (e CanNotReadResponseBodyError) Error() string {
	return "can not read response body: " + e.err.Error()
}
