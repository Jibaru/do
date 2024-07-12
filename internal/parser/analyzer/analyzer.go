package analyzer

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/jibaru/do/internal/types"
)

type Analyzer interface {
	Analyze(expressions types.SectionExpressions) (*types.Sentences, error)
}

type analyzer struct{}

func New() Analyzer {
	return &analyzer{}
}

func (a *analyzer) Analyze(expressions types.SectionExpressions) (*types.Sentences, error) {
	result := types.NewSentences()

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

		if result.Has(key) {
			return nil, NewRepeatedKeyError(key)
		}

		if isStringByQuotes(value) {
			// if value is a string
			result.Set(key, types.String(strings.Trim(value, `"`)))
		} else if isStringByBackticks(value) {
			// if value is a string
			result.Set(key, types.String(strings.Trim(value, "`")))
		} else if isMap(value) {
			// if value is a map
			mp, err := toMap(value)
			if err != nil {
				return nil, err
			}
			result.Set(key, mp)
		} else if isBool(value) {
			// if value is a boolean
			b, _ := strconv.ParseBool(value)
			result.Set(key, types.Bool(b))
		} else if num, ok := isInt(value); ok {
			// if value is an integer
			result.Set(key, types.Int(num))
		} else if floatNum, isOk := isFloat(value); isOk {
			// if value is a number
			result.Set(key, types.Float(floatNum))
		} else if isReferenceToVariable(value) {
			result.Set(key, types.NewReferenceToVariable(value))
		} else if isFunc(value) {
			funcVal, err := toFunc(value)
			if err != nil {
				return nil, err
			}
			result.Set(key, funcVal)
		} else {
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

func isFunc(value string) bool {
	// Regex to match the function call pattern
	re := regexp.MustCompile(`^(\w+)\(([^)]*)\)$`)
	matches := re.FindStringSubmatch(value)

	if matches == nil {
		return false
	}

	if len(matches) < 2 {
		return false
	}

	return matches[1] == types.EnvFuncName ||
		matches[1] == types.FileFuncName ||
		matches[1] == types.UuidFuncName ||
		matches[1] == types.DateFuncName
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
		} else if isFunc(value) {
			funcVal, err := toFunc(val)
			if err != nil {
				return nil, err
			}
			result[key] = funcVal
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

func toFunc(value string) (types.Func, error) {
	re := regexp.MustCompile(`^(\w+)\(([^)]*)\)$`)
	matches := re.FindStringSubmatch(value)

	if matches == nil || len(matches) < 2 {
		return types.Func{}, ReadingExpressionError{parts: matches}
	}

	funcName := matches[1]
	var args []string
	if strings.TrimSpace(matches[2]) != "" {
		// not empty args
		args = strings.Split(matches[2], ",")

		for i := range args {
			args[i] = strings.TrimSpace(args[i])
			args[i] = strings.Trim(args[i], `"`)
		}
	}

	var funcArgs []interface{}
	for _, arg := range args {
		funcArgs = append(funcArgs, types.String(arg))
	}

	return types.NewFunc(funcName, funcArgs)
}
