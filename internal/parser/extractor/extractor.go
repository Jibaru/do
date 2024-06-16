package extractor

import (
	"errors"
	"strings"

	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
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
	sectionNormalizer   normalizer.Normalizer
	partitioner         partitioner.Partitioner
	expressionsAnalyzer analyzer.Analyzer
}

func New(
	sectionNormalizer normalizer.Normalizer,
	partitioner partitioner.Partitioner,
	expressionsAnalyzer analyzer.Analyzer,
) Extractor {
	return &SectionExtractor{
		sectionNormalizer,
		partitioner,
		expressionsAnalyzer,
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

	lines, err := d.partitioner.Split(normalizedContent)
	if err != nil {
		return nil, err
	}

	return d.expressionsAnalyzer.Analyze(lines)
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
