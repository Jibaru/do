package parser

import (
	"errors"
	"strings"
)

var (
	ErrNormalizerEmptyContent = errors.New("empty content")
)

type Normalizer interface {
	Normalize(content string) (string, error)
}

type normalizer struct{}

func NewNormalizer() Normalizer {
	return &normalizer{}
}

func (d *normalizer) Normalize(content string) (string, error) {
	if content = strings.TrimSpace(content); content == "" {
		return "", ErrNormalizerEmptyContent
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

	return result.String(), nil

}
