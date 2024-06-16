package parser

import (
	"github.com/jibaru/do/internal/types"
)

type Mock struct {
	ParseFromFilenameFn func(filename string) (*types.DoFile, error)
}

func (m *Mock) ParseFromFilename(filename string) (*types.DoFile, error) {
	return m.ParseFromFilenameFn(filename)
}
