package caller

type Mock struct {
	CallFn func(variables map[string]interface{}) error
}

func (m *Mock) Call(variables map[string]interface{}) error {
	return m.CallFn(variables)
}
