package reader

import (
	"errors"
	"log"
	"os"
)

var (
	ErrFileReaderCannotRead = errors.New("cannot read file")
)

type FileReader interface {
	Read(filename string) (string, error)
}

type fileReader struct{}

func NewFileReader() FileReader {
	return &fileReader{}
}

func (d *fileReader) Read(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error reading file:", err)
		return "", ErrFileReaderCannotRead
	}

	return string(data), nil
}
