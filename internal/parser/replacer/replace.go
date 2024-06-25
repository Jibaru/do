package replacer

import (
	"fmt"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Replacer interface {
	Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) error
}

type replacer struct{}

func New() Replacer {
	return &replacer{}
}

func (v *replacer) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) error {
	err := v.replaceVariablesInLetSection(letVariables)
	if err != nil {
		return err
	}

	return v.replaceVariablesInDoSection(doVariables, letVariables)
}

func (v *replacer) replaceVariablesInLetSection(letVariables map[string]interface{}) error {
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
				return NewReferenceToVariableNotFoundError(key, val.Value)
			}
			letVariables[key] = variablesWithoutReferences[val.Value]
		}
	}

	return nil
}

func (v *replacer) replaceVariablesInDoSection(doVariables map[string]interface{}, letVariables map[string]interface{}) error {
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
