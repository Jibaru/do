package normalizer

import (
	"github.com/jibaru/do/internal/types"
)

type Mock struct {
	NormalizeFn func(content types.RawSectionContent) (types.NormalizedSectionContent, error)
}

func (m *Mock) Normalize(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
	return m.NormalizeFn(content)
}
