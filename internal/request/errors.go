package request

type CanNotDoRequestError struct {
	err error
}

type CanNotReadResponseBodyError struct {
	err error
}

type CanNotReplaceParamError struct {
	key string
}

func NewCanNotDoRequestError(err error) error {
	return CanNotDoRequestError{err}
}

func NewCanNotReadResponseBodyError(err error) error {
	return CanNotReadResponseBodyError{err}
}

func NewCanNotReplaceParamError(key string) error {
	return CanNotReplaceParamError{key}
}

func (e CanNotDoRequestError) Error() string {
	return "can not do request: " + e.err.Error()
}

func (e CanNotReadResponseBodyError) Error() string {
	return "can not read response body: " + e.err.Error()
}

func (e CanNotReplaceParamError) Error() string {
	return "can not replace param: " + e.key
}
