package analyzer

import "github.com/jibaru/do/internal/types"

type Mock struct {
	AnalyzeFn func(expressions types.SectionExpressions) (*types.Sentences, error)
}

func (m *Mock) Analyze(expressions types.SectionExpressions) (*types.Sentences, error) {
	return m.AnalyzeFn(expressions)
}
