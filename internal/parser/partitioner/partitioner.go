package partitioner

import (
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Partitioner interface {
	Split(content types.NormalizedSectionContent) (types.SectionExpressions, error)
}

type partitioner struct{}

func New() Partitioner {
	return &partitioner{}
}

func (p *partitioner) Split(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
	parts := make(types.SectionExpressions, 0)
	currentPart := strings.Builder{}
	inString := false

	for _, ch := range content {
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

	return parts, nil
}
