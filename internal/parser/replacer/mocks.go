package replacer

type Mock struct {
	ReplaceFn func(doVariables map[string]interface{}, letVariables map[string]interface{})
}

func (m *Mock) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	m.ReplaceFn(doVariables, letVariables)
}
