package reader

import "github.com/jibaru/do/internal/types"

type Mock struct {
	ReadFn func(filename string) (types.FileReaderContent, error)
}

func (m *Mock) Read(filename string) (types.FileReaderContent, error) {
	return m.ReadFn(filename)
}
