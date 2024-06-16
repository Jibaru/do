package normalizer

type Mock struct {
	NormalizeFn func(content string) (string, error)
}

func (m *Mock) Normalize(content string) (string, error) {
	return m.NormalizeFn(content)
}
