package extractor

import (
	"github.com/jibaru/do/internal/types"
)

type Mock struct {
	ExtractFn func(section types.Section, content types.FileReaderContent) (map[string]interface{}, error)
}

func (m *Mock) Extract(section types.Section, content types.FileReaderContent) (map[string]interface{}, error) {
	return m.ExtractFn(section, content)
}
