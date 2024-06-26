package reader

import (
	"os"

	"github.com/jibaru/do/internal/types"
)

type FileReader interface {
	Read(filename string) (types.FileReaderContent, error)
}

type fileReader struct{}

func NewFileReader() FileReader {
	return &fileReader{}
}

func (d *fileReader) Read(filename string) (types.FileReaderContent, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", NewCanNotReadFileError(filename)
	}

	return types.FileReaderContent(data), nil
}
