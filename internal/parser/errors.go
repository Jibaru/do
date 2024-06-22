package parser

type DoSectionEmptyError struct{}
type MethodRequiredError struct{}
type URLRequiredError struct{}
type TypeNotExpectedError struct {
	Key      string
	Expected string
	Actual   string
}

func NewDoSectionEmptyError() error {
	return DoSectionEmptyError{}
}

func NewMethodRequiredError() error {
	return MethodRequiredError{}
}

func NewURLRequiredError() error {
	return URLRequiredError{}
}

func NewTypeNotExpectedError(key, expected, actual string) error {
	return TypeNotExpectedError{
		Key:      key,
		Expected: expected,
		Actual:   actual,
	}
}

func (e DoSectionEmptyError) Error() string {
	return "do section is empty"
}

func (e MethodRequiredError) Error() string {
	return "method is required"
}

func (e URLRequiredError) Error() string {
	return "url is required"
}

func (e TypeNotExpectedError) Error() string {
	return "type not expected for key: " + e.Key + ", expected: " + e.Expected + ", actual: " + e.Actual
}
