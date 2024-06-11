package parser

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

var (
	ErrSectionExtractorNoDoBlock                  = errors.New("no 'do' block found")
	ErrSectionExtractorMissingOpeningBraceAfterDo = errors.New("missing opening brace after 'do'")
	ErrSectionExtractorMissingClosingBrace        = errors.New("missing closing brace for 'do' block")
	ErrSectionExtractorParsingJSON                = errors.New("error parsing JSON value")
	ErrSectionExtractorParsingBooleanValue        = errors.New("error parsing boolean value")
)

type Section string

const (
	LetSection Section = "let"
	DoSection  Section = "do"
)

type SectionExtractor interface {
	Extract(section Section, rawContent string) (map[string]interface{}, error)
}

type sectionExtractor struct {
	doSectionNormalizer Normalizer
}

func NewSectionExtractor(
	doSectionNormalizer Normalizer,
) SectionExtractor {
	return &sectionExtractor{
		doSectionNormalizer,
	}
}

func (d *sectionExtractor) Extract(section Section, rawContent string) (map[string]interface{}, error) {
	content, err := d.extractDoContent(section, rawContent)
	if err != nil {
		return nil, err
	}

	normalizedContent, err := d.doSectionNormalizer.Normalize(content)
	if err != nil {
		return nil, err
	}

	return d.parse(normalizedContent)
}

func (d *sectionExtractor) extractDoContent(section Section, text string) (string, error) {
	startIndex := strings.Index(text, string(section))
	if startIndex == -1 {
		return "", ErrSectionExtractorNoDoBlock
	}

	startIndex += 2 // to skip do
	openBraceIndex := strings.Index(text[startIndex:], "{")
	if openBraceIndex == -1 {
		return "", ErrSectionExtractorMissingOpeningBraceAfterDo
	}
	openBraceIndex += startIndex

	braceCount := 1
	endIndex := openBraceIndex + 1
	content := strings.Builder{}
	for ; endIndex < len(text); endIndex++ {
		if text[endIndex] == '{' {
			braceCount++
		} else if text[endIndex] == '}' {
			braceCount--
			if braceCount == 0 {
				break // found closing brace for do block
			}
		} else {
			content.WriteByte(text[endIndex])
		}
	}

	if endIndex == len(text) {
		return "", ErrSectionExtractorMissingClosingBrace
	}

	return text[openBraceIndex+1 : endIndex], nil
}

func (d *sectionExtractor) parse(normalizedContent string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	lines := strings.Split(normalizedContent, ";")

	for _, line := range lines {
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
