package extractor

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/types"
)

var (
	ErrSectionExtractorNoBlock             = errors.New("no block found")
	ErrSectionExtractorMissingOpeningBrace = errors.New("missing opening brace")
	ErrSectionExtractorMissingClosingBrace = errors.New("missing closing brace")
	ErrSectionExtractorParsingJSON         = errors.New("error parsing JSON value")
	ErrSectionExtractorParsingBooleanValue = errors.New("error parsing boolean value")
)

type Extractor interface {
	Extract(section types.Section, rawContent types.FileReaderContent) (map[string]interface{}, error)
}

type SectionExtractor struct {
	sectionNormalizer normalizer.Normalizer
}

func New(
	sectionNormalizer normalizer.Normalizer,
) Extractor {
	return &SectionExtractor{
		sectionNormalizer,
	}
}

func (d *SectionExtractor) Extract(section types.Section, rawContent types.FileReaderContent) (map[string]interface{}, error) {
	content, err := d.ExtractContent(section, rawContent)
	if err != nil {
		return nil, err
	}

	normalizedContent, err := d.sectionNormalizer.Normalize(content)
	if err != nil {
		if errors.Is(err, normalizer.ErrNormalizerEmptyContent) {
			return nil, ErrSectionExtractorNoBlock
		}
		return nil, err
	}

	return d.parse(normalizedContent)
}

func (d *SectionExtractor) ExtractContent(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
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

	return types.RawSectionContent(sectionContent), nil
}

func (d *SectionExtractor) Parts(normalizedContent types.NormalizedSectionContent) []string {
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

func (d *SectionExtractor) parse(normalizedContent types.NormalizedSectionContent) (map[string]interface{}, error) {
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
