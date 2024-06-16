package replacer

import (
	"fmt"
	"strings"
)

type Replacer interface {
	Replace(doVariables map[string]interface{}, letVariables map[string]interface{})
}

type replacer struct{}

func New() Replacer {
	return &replacer{}
}

func (v *replacer) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	v.replaceVariables(doVariables, letVariables)
}

func (v *replacer) replaceVariables(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	if letVariables == nil {
		return
	}

	for key, value := range doVariables {
		switch val := value.(type) {
		case string:
			doVariables[key] = v.replaceStringVariables(val, letVariables)
		case map[string]interface{}:
			v.replaceVariables(val, letVariables)
		}
	}
}

func (v *replacer) replaceStringVariables(value string, letVariables map[string]interface{}) string {
	for key, val := range letVariables {
		value = strings.ReplaceAll(value, fmt.Sprintf("$%s", key), fmt.Sprintf("%v", val))
	}
	return value
}
