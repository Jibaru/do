package replacer

import "github.com/jibaru/do/internal/types"

type Mock struct {
	ReplaceFn func(doVariables map[string]interface{}, letVariables types.Map) error
}

func (m *Mock) Replace(doVariables map[string]interface{}, letVariables types.Map) error {
	return m.ReplaceFn(doVariables, letVariables)
}
