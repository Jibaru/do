package reader

import (
	"errors"
	"log"
	"os"

	"github.com/jibaru/do/internal/types"
)

var (
	ErrFileReaderCannotRead = errors.New("cannot read file")
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
		log.Println("error reading file:", err)
		return "", ErrFileReaderCannotRead
	}

	return types.FileReaderContent(data), nil
}
