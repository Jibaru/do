package extractor

import (
	"github.com/jibaru/do/internal/types"
)

type Mock struct {
	ExtractFn func(section types.Section, content types.CleanedContent) (*types.Sentences, error)
}

func (m *Mock) Extract(section types.Section, content types.CleanedContent) (*types.Sentences, error) {
	return m.ExtractFn(section, content)
}
