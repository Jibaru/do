package parser

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

var (
	ErrSectionExtractorNoBlock             = errors.New("no block found")
	ErrSectionExtractorMissingOpeningBrace = errors.New("missing opening brace")
	ErrSectionExtractorMissingClosingBrace = errors.New("missing closing brace")
	ErrSectionExtractorParsingJSON         = errors.New("error parsing JSON value")
	ErrSectionExtractorParsingBooleanValue = errors.New("error parsing boolean value")
)

type Section string

const (
	LetSection Section = "let"
	DoSection  Section = "do"
)

type SectionExtractor interface {
	Extract(section Section, rawContent string) (map[string]interface{}, error)
}

type TheSectionExtractor struct {
	doSectionNormalizer Normalizer
}

func NewSectionExtractor(
	doSectionNormalizer Normalizer,
) SectionExtractor {
	return &TheSectionExtractor{
		doSectionNormalizer,
	}
}

func (d *TheSectionExtractor) Extract(section Section, rawContent string) (map[string]interface{}, error) {
	content, err := d.ExtractContent(section, rawContent)
	if err != nil {
		return nil, err
	}

	normalizedContent, err := d.doSectionNormalizer.Normalize(content)
	if err != nil {
		if errors.Is(err, ErrNormalizerEmptyContent) {
			return nil, ErrSectionExtractorNoBlock
		}
		return nil, err
	}

	return d.parse(normalizedContent)
}

func (d *TheSectionExtractor) ExtractContent(section Section, text string) (string, error) {
	inBlock := false
	textInNoBlocks := strings.Builder{}

	for _, ch := range text {
		toWrite := ch
		if inBlock {
			toWrite = ' '
		}

		if ch == '{' {
			inBlock = true
			toWrite = ' '
		}

		if ch == '}' {
			inBlock = false
			toWrite = ' '
		}

		textInNoBlocks.WriteRune(toWrite)
	}

	textWithOnlyBlocks := textInNoBlocks.String()
	startIndex := strings.Index(textWithOnlyBlocks, string(section))
	if startIndex == -1 {
		return "", ErrSectionExtractorNoBlock
	}

	content := strings.Builder{}
	inBlock = false
	inString := false
	openBracesCount := 0
	foundOpeningBrace := false

	for i := startIndex + len(string(section)); i < len(text); i++ {
		if text[i] == '{' && !inString {
			if !foundOpeningBrace {
				foundOpeningBrace = true
				inBlock = true
			}
			openBracesCount++
			if openBracesCount > 1 {
				content.WriteByte(text[i])
			}
			continue
		}

		if text[i] == '}' && !inString {
			openBracesCount--
			if openBracesCount == 0 {
				inBlock = false
				break
			}
			content.WriteByte(text[i])
			continue
		}

		if text[i] == '"' {
			inString = !inString
		}

		if inBlock {
			content.WriteByte(text[i])
		}
	}

	if !foundOpeningBrace {
		return "", ErrSectionExtractorMissingOpeningBrace
	}

	if openBracesCount != 0 {
		return "", ErrSectionExtractorMissingClosingBrace
	}

	sectionContent := content.String()

	if strings.TrimSpace(sectionContent) == "" {
		return "", ErrSectionExtractorNoBlock
	}

	return sectionContent, nil
}

func (d *TheSectionExtractor) Parts(normalizedContent string) []string {
	parts := make([]string, 0)
	currentPart := strings.Builder{}
	inString := false

	for _, ch := range normalizedContent {
		if ch == '"' {
			inString = !inString
		}

		if ch == ';' && !inString {
			part := currentPart.String()
			currentPart.Reset()

			if part == ";" {
				continue
			}

			parts = append(parts, part)
			continue
		}

		currentPart.WriteRune(ch)
	}

	return parts
}

func (d *TheSectionExtractor) parse(normalizedContent string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	lines := d.Parts(normalizedContent)

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
