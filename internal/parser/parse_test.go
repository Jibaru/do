package parser_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

func TestParser_FromFilename(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expected      *types.DoFile
		expectedError error
		FileReaderFn  func(filename string) (string, error)
		ExtractorFn   func(section extractor.Section, content string) (map[string]interface{}, error)
		ReplacerFn    func(doVariables map[string]interface{}, letVariables map[string]interface{})
	}{
		{
			name:     "success",
			filename: "valid.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: map[string]interface{}{
						"var1": 12,
						"var2": "text",
						"var3": false,
						"var4": 12.33,
					},
				},
				Do: types.Do{
					Method:  "GET",
					URL:     "http://localhost:8080/api/todos/:id",
					Params:  map[string]interface{}{"id": "12"},
					Query:   map[string]interface{}{"isOk": "false"},
					Headers: map[string]interface{}{"Authorization": "Bearer text"},
					Body:    "{\"extra\": 12, \"extra2\": false, \"extra3\": \"text\", \"extra4\": 12.33}",
				},
			},
			FileReaderFn: func(filename string) (string, error) {
				return "let{var1=12;var2=\"text\";var3=false;var4=12.33;}do{method=\"GET\";url=\"http://localhost:8080/api/todos/:id\";params={\"id\":\"$id\"};query={\"isOk\":\"$isOk\"};headers={\"Authorization\":\"Bearer $token\"};body=`{\"extra\": $extra, \"extra2\": $extra2, \"extra3\": \"$extra3\", \"extra4\": $extra4}`;}", nil
			},
			ExtractorFn: func(section extractor.Section, content string) (map[string]interface{}, error) {
				if section == extractor.LetSection {
					return map[string]interface{}{
						"var1": 12,
						"var2": "text",
						"var3": false,
						"var4": 12.33,
					}, nil
				}

				if section == extractor.DoSection {
					return map[string]interface{}{
						"method":  "GET",
						"url":     "http://localhost:8080/api/todos/:id",
						"params":  map[string]interface{}{"id": "$id"},
						"query":   map[string]interface{}{"isOk": "$isOk"},
						"headers": map[string]interface{}{"Authorization": "Bearer $token"},
						"body":    "{\"extra\": $var1, \"extra2\": $var2, \"extra3\": \"$var3\", \"extra4\": $var4}",
					}, nil
				}

				return nil, nil
			},
			ReplacerFn: func(doVariables map[string]interface{}, letVariables map[string]interface{}) {
				doVariables["params"] = map[string]interface{}{"id": "12"}
				doVariables["query"] = map[string]interface{}{"isOk": "false"}
				doVariables["headers"] = map[string]interface{}{"Authorization": "Bearer text"}
				doVariables["body"] = "{\"extra\": 12, \"extra2\": false, \"extra3\": \"text\", \"extra4\": 12.33}"
			},
		},
	}

	fileReader := &reader.MockFileReader{}
	sectionExtractor := &extractor.Mock{}
	varReplacer := &replacer.Mock{}

	p := parser.New(fileReader, sectionExtractor, varReplacer)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileReader.ReadFn = tc.FileReaderFn
			sectionExtractor.ExtractFn = tc.ExtractorFn
			varReplacer.ReplaceFn = tc.ReplacerFn

			doFile, err := p.FromFilename(tc.filename)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if doFile != nil && !reflect.DeepEqual(doFile, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, doFile)
			}
		})
	}
}
