package parser

type DoSectionEmptyError struct{}
type MethodRequiredError struct{}
type URLRequiredError struct{}

func NewDoSectionEmptyError() error {
	return DoSectionEmptyError{}
}

func NewMethodRequiredError() error {
	return MethodRequiredError{}
}

func NewURLRequiredError() error {
	return URLRequiredError{}
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
