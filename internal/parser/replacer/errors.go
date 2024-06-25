package replacer

type ReferenceToVariableNotFoundError struct {
	key           string
	referenceName string
}

func NewReferenceToVariableNotFoundError(key, referenceName string) error {
	return ReferenceToVariableNotFoundError{key, referenceName}
}

func (e ReferenceToVariableNotFoundError) Error() string {
	return "reference to variable for key " + e.key + " not found: " + e.referenceName
}
