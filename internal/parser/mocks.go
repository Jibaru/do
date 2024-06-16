package parser

import (
	"github.com/jibaru/do/internal/types"
)

type Mock struct {
	FromFilenameFn func(filename string) (*types.DoFile, error)
}

func (m *Mock) FromFilename(filename string) (*types.DoFile, error) {
	return m.FromFilenameFn(filename)
}
