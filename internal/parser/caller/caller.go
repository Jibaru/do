package caller

import (
	"errors"

	"github.com/jibaru/do/internal/types"
)

type Caller interface {
	Call(variables map[string]interface{}) error
}

type caller struct {
}

func New() Caller {
	return &caller{}
}

func (c *caller) Call(variables map[string]interface{}) error {
	for name, value := range variables {
		switch value.(type) {
		case types.Func:
			fn := value.(types.Func)

			if fn.HasReferences() {
				return errors.New("function for key " + name + " has references")
			}

			resolvedFn, err := fn.Resolve()
			if err != nil {
				return err
			}

			switch resolvedFn.(type) {
			case types.EnvFunc:
				envFunc := resolvedFn.(types.EnvFunc)
				variables[name] = envFunc.Exec()
			case types.FileFunc:
				fileFunc := resolvedFn.(types.FileFunc)
				variables[name] = fileFunc.Exec()
			}
		}
	}

	return nil
}
