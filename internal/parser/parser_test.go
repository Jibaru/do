package parser_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

func TestParser_FromFilename(t *testing.T) {
	body := types.String("{\"extra\": 12, \"extra2\": false, \"extra3\": \"text\", \"extra4\": 12.33}")

	testCases := []struct {
		name          string
		filename      string
		expected      *types.DoFile
		expectedError error
		FileReaderFn  func(filename string) (types.FileReaderContent, error)
		ExtractorFn   func(section types.Section, content types.FileReaderContent) (map[string]interface{}, error)
		ReplacerFn    func(doVariables map[string]interface{}, letVariables map[string]interface{}) error
	}{
		{
			name:     "success",
			filename: "valid.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: map[string]interface{}{
						"var1": types.Int(12),
						"var2": types.String("text"),
						"var3": types.Bool(false),
						"var4": types.Float(12.33),
					},
				},
				Do: types.Do{
					Method:  types.String("GET"),
					URL:     types.String("http://localhost:8080/api/todos/:id"),
					Params:  types.Map{"id": types.String("12")},
					Query:   types.Map{"isOk": types.String("false")},
					Headers: types.Map{"Authorization": types.String("Bearer text")},
					Body:    utils.Ptr(body),
				},
			},
			FileReaderFn: func(filename string) (types.FileReaderContent, error) {
				return "let{var1=12;var2=\"text\";var3=false;var4=12.33;}do{method=\"GET\";url=\"http://localhost:8080/api/todos/:id\";params={\"id\":\"$id\"};query={\"isOk\":\"$isOk\"};headers={\"Authorization\":\"Bearer $token\"};body=`{\"extra\": $extra, \"extra2\": $extra2, \"extra3\": \"$extra3\", \"extra4\": $extra4}`;}", nil
			},
			ExtractorFn: func(section types.Section, content types.FileReaderContent) (map[string]interface{}, error) {
				if section == types.LetSection {
					return map[string]interface{}{
						"var1": types.Int(12),
						"var2": types.String("text"),
						"var3": types.Bool(false),
						"var4": types.Float(12.33),
					}, nil
				}

				if section == types.DoSection {
					return map[string]interface{}{
						"method":  types.String("GET"),
						"url":     types.String("http://localhost:8080/api/todos/:id"),
						"params":  types.Map{"id": types.String("$id")},
						"query":   types.Map{"isOk": types.String("$isOk")},
						"headers": types.Map{"Authorization": types.String("Bearer $token")},
						"body":    types.String("{\"extra\": $var1, \"extra2\": $var2, \"extra3\": \"$var3\", \"extra4\": $var4}"),
					}, nil
				}

				return nil, nil
			},
			ReplacerFn: func(doVariables map[string]interface{}, letVariables map[string]interface{}) error {
				doVariables["params"] = types.Map{"id": types.String("12")}
				doVariables["query"] = types.Map{"isOk": types.String("false")}
				doVariables["headers"] = types.Map{"Authorization": types.String("Bearer text")}
				doVariables["body"] = body
				return nil
			},
		},
	}

	fileReader := &reader.Mock{}
	sectionExtractor := &extractor.Mock{}
	varReplacer := &replacer.Mock{}

	p := parser.New(fileReader, sectionExtractor, varReplacer)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileReader.ReadFn = tc.FileReaderFn
			sectionExtractor.ExtractFn = tc.ExtractorFn
			varReplacer.ReplaceFn = tc.ReplacerFn

			doFile, err := p.ParseFromFilename(tc.filename)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if doFile != nil && tc.expected == nil {
				t.Errorf("expected %v, got %v", tc.expected, doFile)
			} else if doFile == nil && tc.expected != nil {
				t.Errorf("expected %v, got %v", tc.expected, doFile)
			} else if doFile != nil && tc.expected != nil {
				if !reflect.DeepEqual(doFile, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, doFile)
				}
			}
		})
	}
}
