package taker

import "github.com/jibaru/do/internal/types"

type Mock struct {
	TakeFn func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error)
}

func (m *Mock) Take(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
	return m.TakeFn(section, text)
}
