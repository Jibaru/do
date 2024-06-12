package parser

type MockNormalizer struct {
	NormalizeFn func(content string) (string, error)
}

func (m *MockNormalizer) Normalize(content string) (string, error) {
	return m.NormalizeFn(content)
}
