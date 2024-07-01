package replacer

import (
	"fmt"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type DoReplacer interface {
	// Replace replaces the variables in the do section from the let section
	Replace(doVariables map[string]interface{}, letVariables types.Map) error
}

type replacer struct{}

func New() DoReplacer {
	return &replacer{}
}

func (v *replacer) Replace(doVariables map[string]interface{}, letVariables types.Map) error {
	if letVariables != nil && !letVariables.HasBasicTypesValues() {
		return NewInvalidLetVariablesError()
	}

	return v.replaceVariablesInDoSection(doVariables, letVariables)
}

func (v *replacer) replaceVariablesInDoSection(doVariables map[string]interface{}, letVariables types.Map) error {
	if letVariables == nil {
		return nil
	}

	for key, value := range doVariables {
		switch val := value.(type) {
		case types.String:
			doVariables[key] = v.replaceStringVariables(val, letVariables)
		case types.ReferenceToVariable:
			if _, ok := letVariables[val.Value]; !ok {
				return NewReferenceToVariableNotFoundError(key, val.Value)
			}
			doVariables[key] = letVariables[val.Value]
		case types.Map:
			err := v.replaceVariablesInDoSection(val, letVariables)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
