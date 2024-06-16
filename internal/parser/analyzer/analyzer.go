package analyzer

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/jibaru/do/internal/types"
)

var (
	ErrParsingMapValue     = errors.New("error parsing JSON value")
	ErrParsingBooleanValue = errors.New("error parsing boolean value")
)

type Analyzer interface {
	Analyze(expressions types.SectionExpressions) (map[string]interface{}, error)
}

type analyzer struct{}

func New() Analyzer {
	return &analyzer{}
}

func isString(value string) bool {
	return (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
		(strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`"))
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

		if isString(value) {
			// if value is a string
			result[key] = strings.Trim(value, `"`)
		} else if isMap(value) {
			// if value is a JSON object
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(value), &obj); err != nil {
				return nil, NewCanNotParseMapValueError(err)
			}
			result[key] = obj
		} else if isBool(value) {
			// if value is a boolean
			b, err := strconv.ParseBool(value)
			if err != nil {
				return nil, NewCanNotParseBoolValueError(err)
			}
			result[key] = b
		} else if num, ok := isInt(value); ok {
			// if value is an integer
			result[key] = int(num)
		} else if floatNum, isOk := isFloat(value); isOk {
			// if value is a number
			result[key] = floatNum
		} else {
			// otherwise, value is a string
			result[key] = value
		}
	}

	return result, nil
}
