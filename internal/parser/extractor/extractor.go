package extractor

import (
	"errors"

	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/types"
)

var (
	ErrSectionExtractorNoBlock = errors.New("no block found")
)

type Extractor interface {
	Extract(section types.Section, rawContent types.FileReaderContent) (map[string]interface{}, error)
}

type SectionExtractor struct {
	sectionTaker        taker.Taker
	sectionNormalizer   normalizer.Normalizer
	partitioner         partitioner.Partitioner
	expressionsAnalyzer analyzer.Analyzer
}

func New(
	sectionTaker taker.Taker,
	sectionNormalizer normalizer.Normalizer,
	partitioner partitioner.Partitioner,
	expressionsAnalyzer analyzer.Analyzer,
) Extractor {
	return &SectionExtractor{
		sectionTaker,
		sectionNormalizer,
		partitioner,
		expressionsAnalyzer,
	}
}

func (d *SectionExtractor) Extract(section types.Section, rawContent types.FileReaderContent) (map[string]interface{}, error) {
	content, err := d.sectionTaker.Take(section, rawContent)
	if err != nil && errors.Is(err, taker.NoBlockError{}) {
		return nil, ErrSectionExtractorNoBlock
	} else if err != nil {
		return nil, err
	}

	normalizedContent, err := d.sectionNormalizer.Normalize(content)
	if err != nil {
		if errors.Is(err, normalizer.EmptyContentError{}) {
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
