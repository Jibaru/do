package resolver

type InvalidVariablesError struct {
	reason string
}

type ReferenceToVariableNotFoundError struct {
	key   string
	value string
}

func NewInvalidVariablesError(reason string) InvalidVariablesError {
	return InvalidVariablesError{reason}
}

func NewReferenceToVariableNotFoundError(key, value string) ReferenceToVariableNotFoundError {
	return ReferenceToVariableNotFoundError{key, value}
}

func (e InvalidVariablesError) Error() string {
	return "invalid variables error: " + e.reason
}

func (e ReferenceToVariableNotFoundError) Error() string {
	return "reference to variable not found error: " + e.key + ", variable: " + e.value
}
