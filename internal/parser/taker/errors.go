package taker

type NoBlockError struct{}
type MissingOpeningBraceError struct{}
type MissingClosingBraceError struct{}

func NewNoBlockError() error {
	return NoBlockError{}
}

func NewMissingOpeningBraceError() error {
	return MissingOpeningBraceError{}
}

func NewMissingClosingBraceError() error {
	return MissingClosingBraceError{}
}

func (e NoBlockError) Error() string {
	return "no block found"
}

func (e MissingOpeningBraceError) Error() string {
	return "missing opening brace"
}

func (e MissingClosingBraceError) Error() string {
	return "missing closing brace"
}
