package parser

import "github.com/jibaru/do/internal/types"

type MockNormalizer struct {
	NormalizeFn func(content string) (string, error)
}

func (m *MockNormalizer) Normalize(content string) (string, error) {
	return m.NormalizeFn(content)
}

type MockSectionExtractor struct {
	ExtractFn func(section Section, content string) (map[string]interface{}, error)
}

func (m *MockSectionExtractor) Extract(section Section, content string) (map[string]interface{}, error) {
	return m.ExtractFn(section, content)
}

type MockVariablesReplacer struct {
	ReplaceFn func(doVariables map[string]interface{}, letVariables map[string]interface{})
}

func (m *MockVariablesReplacer) Replace(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	m.ReplaceFn(doVariables, letVariables)
}

type MockParser struct {
	FromFilenameFn func(filename string) (*types.DoFile, error)
}

func (m *MockParser) FromFilename(filename string) (*types.DoFile, error) {
	return m.FromFilenameFn(filename)
}
