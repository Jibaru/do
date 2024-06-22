package analyzer

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Analyzer interface {
	Analyze(expressions types.SectionExpressions) (map[string]interface{}, error)
}

type analyzer struct{}

func New() Analyzer {
	return &analyzer{}
}

func isStringByQuotes(value string) bool {
	return (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`))
}

func isStringByBackticks(value string) bool {
	return (strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`"))
}

func isMap(value string) bool {
	return strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")
}

func isBool(value string) bool {
	return value == "true" || value == "false"
}

func isInt(value string) (int, bool) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return int(val), true
	}

	return 0, false
}

func isFloat(value string) (float64, bool) {
	val, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return val, true
	}

	return 0, false
}

func (a *analyzer) Analyze(expressions types.SectionExpressions) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, line := range expressions {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, NewReadingExpressionError(parts...)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, ok := result[key]; ok {
			return nil, NewRepeatedKeyError(key)
		}

		if isStringByQuotes(value) {
			// if value is a string
			result[key] = strings.Trim(value, `"`)
		} else if isStringByBackticks(value) {
			// if value is a string
			result[key] = strings.Trim(value, "`")
		} else if isMap(value) {
			// if value is a JSON object
			var obj map[string]interface{}
			_ = json.Unmarshal([]byte(value), &obj)
			result[key] = obj
		} else if isBool(value) {
			// if value is a boolean
			b, _ := strconv.ParseBool(value)
			result[key] = b
		} else if num, ok := isInt(value); ok {
			// if value is an integer
			result[key] = int(num)
		} else if floatNum, isOk := isFloat(value); isOk {
			// if value is a number
			result[key] = floatNum
		} else {
			// otherwise
			return nil, NewInvalidValueError(value)
		}
	}

	return result, nil
}
