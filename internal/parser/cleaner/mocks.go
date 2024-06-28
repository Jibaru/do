package cleaner

import "github.com/jibaru/do/internal/types"

type Mock struct {
	CleanFn func(rawContent types.FileReaderContent) (types.CleanedContent, error)
}

func (m *Mock) Clean(rawContent types.FileReaderContent) (types.CleanedContent, error) {
	return m.CleanFn(rawContent)
}
