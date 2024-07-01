package resolver

import "github.com/jibaru/do/internal/types"

type Mock struct {
	ResolveFn func(sentences *types.Sentences) (*types.Sentences, error)
}

func (m *Mock) Resolve(sentences *types.Sentences) (*types.Sentences, error) {
	return m.ResolveFn(sentences)
}
