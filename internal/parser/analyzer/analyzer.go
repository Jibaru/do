package analyzer

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/jibaru/do/internal/types"
)

var (
	ErrSectionExtractorParsingJSON         = errors.New("error parsing JSON value")
	ErrSectionExtractorParsingBooleanValue = errors.New("error parsing boolean value")
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
			log.Println("error reading do section line:", parts, len(parts))
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			// if value is a string
			result[key] = strings.Trim(value, `"`)
		} else if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
			// if value is a string
			result[key] = strings.Trim(value, "`")
		} else if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			// if value is a JSON object
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(value), &obj); err != nil {
				return nil, ErrSectionExtractorParsingJSON
			}
			result[key] = obj
		} else if value == "true" || value == "false" {
			// if value is a boolean
			b, err := strconv.ParseBool(value)
			if err != nil {
				return nil, ErrSectionExtractorParsingBooleanValue
			}
			result[key] = b
		} else if num, err := strconv.ParseInt(value, 10, 64); err == nil {
			// if value is an integer
			result[key] = int(num)
		} else if num, err := strconv.ParseFloat(value, 64); err == nil {
			// if value is a number
			result[key] = num
		} else {
			// otherwise, value is a string
			result[key] = value
		}
	}

	return result, nil
}
