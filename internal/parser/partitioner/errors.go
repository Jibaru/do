package partitioner

type EmptyPartError struct{}

func NewEmptyPartFound() error {
	return EmptyPartError{}
}

func (e EmptyPartError) Error() string {
	return "empty part"
}
