package partitioner

type EmptyPartFound struct{}

func NewEmptyPartFound() error {
	return EmptyPartFound{}
}

func (e EmptyPartFound) Error() string {
	return "empty part found"
}
