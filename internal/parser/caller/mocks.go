package caller

type Mock struct {
	CallFn func(letVariables map[string]interface{}, doVariables map[string]interface{}) error
}

func (m *Mock) Call(letVariables map[string]interface{}, doVariables map[string]interface{}) error {
	return m.CallFn(letVariables, doVariables)
}
