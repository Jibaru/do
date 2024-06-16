package extractor

type Mock struct {
	ExtractFn func(section Section, content string) (map[string]interface{}, error)
}

func (m *Mock) Extract(section Section, content string) (map[string]interface{}, error) {
	return m.ExtractFn(section, content)
}
