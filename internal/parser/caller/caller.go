package caller

import "github.com/jibaru/do/internal/types"

type Caller interface {
	Call(letVariables map[string]interface{}, doVariables map[string]interface{}) error
}

type caller struct {
}

func New() Caller {
	return &caller{}
}

func (c *caller) Call(letVariables map[string]interface{}, doVariables map[string]interface{}) error {
	for name, value := range letVariables {
		switch value.(type) {
		case types.EnvFunc:
			envFunc := value.(types.EnvFunc)
			letVariables[name] = types.String(envFunc.Exec())
		}
	}

	for name, value := range doVariables {
		switch value.(type) {
		case types.EnvFunc:
			envFunc := value.(types.EnvFunc)
			doVariables[name] = types.String(envFunc.Exec())
		}
	}

	return nil
}
