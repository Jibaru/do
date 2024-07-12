package caller

import (
	"errors"

	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

type Caller interface {
	Call(variables map[string]interface{}) error
}

type caller struct {
	uuidFactory utils.UuidFactory
	dateFactory utils.DateFactory
}

func New(
	uuidFactory utils.UuidFactory,
	dateFactory utils.DateFactory,
) Caller {
	return &caller{
		uuidFactory,
		dateFactory,
	}
}

func (c *caller) Call(variables map[string]interface{}) error {
	for name, value := range variables {
		switch value.(type) {
		case types.Func:
			fn := value.(types.Func)

			if fn.HasReferences() {
				return errors.New("function for key " + name + " has references")
			}

			resolvedFn, err := fn.Resolve(
				c.uuidFactory,
				c.dateFactory,
			)
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
			case types.UuidFunc:
				uuidFunc := resolvedFn.(types.UuidFunc)
				variables[name] = uuidFunc.Exec()
			case types.DateFunc:
				dateFunc := resolvedFn.(types.DateFunc)
				variables[name] = dateFunc.Exec()
			}
		}
	}

	return nil
}
