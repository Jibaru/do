package normalizer

import (
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Normalizer interface {
	Normalize(content types.RawSectionContent) (types.NormalizedSectionContent, error)
}

type normalizer struct{}

func New() Normalizer {
	return &normalizer{}
}

func (d *normalizer) Normalize(rawSectionContent types.RawSectionContent) (types.NormalizedSectionContent, error) {
	content := string(rawSectionContent)

	if content = strings.TrimSpace(content); content == "" {
		return "", NewEmptyContentError()
	}

	var result strings.Builder
	inQuotes := false

	for i, ch := range content {
		if ch == '"' {
			inQuotes = !inQuotes
		}
		if inQuotes {
			result.WriteByte(content[i])
			continue
		}

		if content[i] != ' ' && content[i] != '\t' && content[i] != '\n' && content[i] != '\r' {
			result.WriteByte(content[i])
		}
	}

	return types.NormalizedSectionContent(result.String()), nil
}
