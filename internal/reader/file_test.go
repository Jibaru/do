package reader_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

func TestFileReader_Read(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expected      types.FileReaderContent
		expectedError error
	}{
		{
			name:     "success",
			filename: "testdata/01.do",
			expected: "let{var1=12;var2=\"text\";var3=false;var4=12.33;}do{method=\"GET\";url=\"http://localhost:8080/api/todos/:id\";params={\"id\":\"$id\"};query={\"id\":\"$id\"};headers={\"Authorization\":\"$token\",\"Content-Type\":\"application/json\"};body=`{\"extra\":\"something\"}`;}",
		},
		{
			name:          "error file not found",
			filename:      "testdata/not_found.do",
			expectedError: errors.New("can not read file: open testdata/not_found.do: no such file or directory"),
		},
	}

	fileReader := reader.NewFileReader()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := fileReader.Read(tc.filename)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if content != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, content)
			}
		})
	}
}
