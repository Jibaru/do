package replacer

import (
	"fmt"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Replacer interface {
	Replace(doVariables map[string]interface{}, letVariables map[string]interface{})
}

type replacer struct{}

func New() Replacer {
	return &replacer{}
}

func (v *replacer) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	variablesWithoutReferences := make(map[string]interface{})

	for key, value := range letVariables {
		switch val := value.(type) {
		case types.ReferenceToVariable:
			continue
		default:
			variablesWithoutReferences[key] = val
		}
	}

	for key, value := range letVariables {
		switch val := value.(type) {
		case types.ReferenceToVariable:
			if _, ok := variablesWithoutReferences[val.Value]; !ok {
				// TODO: raise error variable not found
				continue
			}
			letVariables[key] = variablesWithoutReferences[val.Value]
		}
	}

	v.replaceVariables(doVariables, letVariables)
}

func (v *replacer) replaceVariables(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	if letVariables == nil {
		return
	}

	for key, value := range doVariables {
		switch val := value.(type) {
		case types.String:
			doVariables[key] = v.replaceStringVariables(val, letVariables)
		case types.ReferenceToVariable:
			if _, ok := letVariables[val.Value]; !ok {
				// TODO: raise error variable not found
				continue
			}
			doVariables[key] = letVariables[val.Value]
		case types.Map:
			v.replaceVariables(val, letVariables)
		}
	}
}

func (v *replacer) replaceStringVariables(value types.String, letVariables map[string]interface{}) types.String {
	for key, val := range letVariables {
		stringVal := fmt.Sprintf("%v", val)

		value = types.String(
			strings.ReplaceAll(
				string(value),
				fmt.Sprintf("$%s", key),
				fmt.Sprintf("%v", stringVal),
			),
		)
	}
	return value
}
