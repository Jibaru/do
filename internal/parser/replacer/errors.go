package replacer

type ReferenceToVariableNotFoundError struct {
	key           string
	referenceName string
}

type LetVariablesNotBasicTypesError struct{}

func NewReferenceToVariableNotFoundError(key, referenceName string) error {
	return ReferenceToVariableNotFoundError{key, referenceName}
}

func (e ReferenceToVariableNotFoundError) Error() string {
	return "reference to variable for key " + e.key + " not found: " + e.referenceName
}

func NewInvalidLetVariablesError() error {
	return LetVariablesNotBasicTypesError{}
}

func (e LetVariablesNotBasicTypesError) Error() string {
	return "let variables must have basic types values"
}
