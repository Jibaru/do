package analyzer

import "github.com/jibaru/do/internal/types"

type Mock struct {
	AnalyzeFn func(expressions types.SectionExpressions) (map[string]interface{}, error)
}

func (m *Mock) Analyze(expressions types.SectionExpressions) (map[string]interface{}, error) {
	return m.AnalyzeFn(expressions)
}
