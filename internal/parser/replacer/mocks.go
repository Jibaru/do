package replacer

type Mock struct {
	ReplaceFn func(doVariables map[string]interface{}, letVariables map[string]interface{}) error
}

func (m *Mock) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) error {
	return m.ReplaceFn(doVariables, letVariables)
}
