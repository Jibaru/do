package parser

import (
	"errors"
	"fmt"
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

	return d.normalizeContent(content)
}

func (d *normalizer) normalizeContent(doContent string) (string, error) {
	var result strings.Builder
	lines := strings.Split(strings.TrimSpace(doContent), ";")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = d.removeSpacesExceptInQuotes(line)
		result.WriteString(fmt.Sprintf("%s;", line))
	}

	return result.String(), nil
}

func (d *normalizer) removeSpacesExceptInQuotes(line string) string {
	var result strings.Builder
	var insideQuotes bool = false
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			insideQuotes = !insideQuotes
		}

		if insideQuotes {
			result.WriteByte(line[i])
			continue
		}

		if line[i] != ' ' && line[i] != '\t' && line[i] != '\n' && line[i] != '\r' {
			result.WriteByte(line[i])
		}
	}
	return result.String()
}
