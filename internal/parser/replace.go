package parser

import (
	"fmt"
	"strings"
)

type VariablesReplacer interface {
	Replace(doVariables map[string]interface{}, letVariables map[string]interface{})
}

type variablesReplacer struct{}

func NewVariablesReplacer() VariablesReplacer {
	return &variablesReplacer{}
}

func (v *variablesReplacer) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	replaceVariables(doVariables, letVariables)
}

func replaceVariables(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	if letVariables == nil {
		return
	}

	for key, value := range doVariables {
		switch val := value.(type) {
		case string:
			doVariables[key] = replaceStringVariables(val, letVariables)
		case map[string]interface{}:
			replaceVariables(val, letVariables)
		}
	}
}

func replaceStringVariables(value string, letVariables map[string]interface{}) string {
	for key, val := range letVariables {
		value = strings.ReplaceAll(value, fmt.Sprintf("$%s", key), fmt.Sprintf("%v", val))
	}
	return value
}
