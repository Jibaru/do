package partitioner

import "github.com/jibaru/do/internal/types"

type Mock struct {
	SplitFn func(content types.NormalizedSectionContent) (types.SectionExpressions, error)
}

func (m *Mock) Split(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
	return m.SplitFn(content)
}
