package analyzer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/jibaru/do/internal/types"
)

type Analyzer interface {
	Analyze(expressions types.SectionExpressions) (map[string]interface{}, error)
}

type analyzer struct{}

func New() Analyzer {
	return &analyzer{}
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

		if types.IsReservedKeyword(key) {
			return nil, NewReservedKeywordError(key)
		}

		if _, ok := result[key]; ok {
			return nil, NewRepeatedKeyError(key)
		}

		if isStringByQuotes(value) {
			// if value is a string
			result[key] = types.String(strings.Trim(value, `"`))
		} else if isStringByBackticks(value) {
			// if value is a string
			result[key] = types.String(strings.Trim(value, "`"))
		} else if isMap(value) {
			// if value is a map
			mp, err := toMap(value)
			if err != nil {
				return nil, err
			}
			result[key] = mp
		} else if isBool(value) {
			// if value is a boolean
			b, _ := strconv.ParseBool(value)
			result[key] = types.Bool(b)
		} else if num, ok := isInt(value); ok {
			// if value is an integer
			result[key] = types.Int(num)
		} else if floatNum, isOk := isFloat(value); isOk {
			// if value is a number
			result[key] = types.Float(floatNum)
		} else if isReferenceToVariable(value) {
			result[key] = types.NewReferenceToVariable(value)
		} else {
			// otherwise
			return nil, NewInvalidValueError(value)
		}
	}

	return result, nil
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

func toMap(value string) (types.Map, error) {
	value = strings.TrimSpace(value[1 : len(value)-1]) // Remove the curly braces
	parts := strings.Split(value, ",")
	result := make(map[string]interface{})

	for _, part := range parts {
		pair := strings.SplitN(part, ":", 2)
		if len(pair) != 2 {
			return nil, NewInvalidValueError(part)
		}

		key := strings.TrimSpace(pair[0])
		key = strings.Trim(key, `"`)
		val := strings.TrimSpace(pair[1])

		if isStringByQuotes(val) {
			result[key] = types.String(strings.Trim(val, `"`))
		} else if isStringByBackticks(val) {
			result[key] = types.String(strings.Trim(val, "`"))
		} else if isMap(val) {
			obj, err := toMap(val)
			if err != nil {
				return nil, err
			}
			result[key] = obj
		} else if isBool(val) {
			b, _ := strconv.ParseBool(val)
			result[key] = types.Bool(b)
		} else if num, ok := isInt(val); ok {
			result[key] = types.Int(num)
		} else if floatNum, isOk := isFloat(val); isOk {
			result[key] = types.Float(floatNum)
		} else if isReferenceToVariable(val) {
			result[key] = types.NewReferenceToVariable(val)
		} else {
			return nil, NewInvalidValueError(val)
		}
	}

	return types.Map(result), nil
}

func isReferenceToVariable(value string) bool {
	for i, char := range value {
		if i == 0 && unicode.IsDigit(char) {
			return false
		}

		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_') {
			return false
		}
	}

	return true
}
