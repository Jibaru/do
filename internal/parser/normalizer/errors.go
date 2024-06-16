package normalizer

type EmptyContentError struct{}

func NewEmptyContentError() error {
	return EmptyContentError{}
}

func (e EmptyContentError) Error() string {
	return "empty content"
}
