package parser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

func TestParser_ParseFromFilename_Integration(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		setup    func(t *testing.T)
		expected *types.DoFile
	}{
		{
			name: "01_no_variables.do",
			path: "examples/01_no_variables.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: nil,
				},
				Do: types.Do{
					Method: types.String("GET"),
					URL:    types.String("https://jsonplaceholder.typicode.com/todos/:id"),
					Params: map[string]interface{}{
						"id": types.String("1"),
					},
				},
			},
		},
		{
			name: "02_multiple_variables.do",
			path: "examples/02_multiple_variables.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: map[string]interface{}{
						"var1": types.Int(1),
						"var2": types.String("hello"),
						"var3": types.Bool(true),
						"var4": types.Bool(false),
						"var5": types.ReferenceToVariable{Value: "var1"},
					},
				},
				Do: types.Do{
					Method: types.String("GET"),
					URL:    types.String("http://example.com/:id"),
					Params: map[string]interface{}{
						"id": types.String("1"),
					},
					Headers: map[string]interface{}{
						"Content-Type": types.String("application/json"),
						"X-Message":    types.String("hello"),
						"X-Var5":       types.Int(1),
					},
					Body: utils.Ptr(types.String(`{"var1":1,"var2":"hello","var3":true,"var4":false,"var5":1}`)),
				},
			},
		},
		{
			name: "03_braces.do",
			path: "examples/03_braces.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: map[string]interface{}{
						"var1": types.String("{};"),
						"var2": types.String("{;;;"),
					},
				},
				Do: types.Do{
					Method: types.String("GET"),
					URL:    types.String("http://localhost:8080/{};"),
				},
			},
		},
		{
			name: "04_functions.do",
			path: "examples/04_functions.do",
			setup: func(t *testing.T) {
				err := os.Setenv("ALREADY_EXIST", "already_exist")
				if err != nil {

				}
			},
			expected: &types.DoFile{
				Let: types.Let{
					Variables: map[string]interface{}{
						"var1": types.String("default"),
						"path": types.String("already_exist"),
					},
				},
				Do: types.Do{
					Method: types.String("GET"),
					URL:    types.String("https://jsonplaceholder.typicode.com/todos/:id"),
					Query: map[string]interface{}{
						"id":  types.String("default"),
						"id2": types.String("default"),
					},
				},
			},
		},
	}

	doFileReader := reader.NewFileReader()
	commentCleaner := cleaner.New()
	sectionTaker := taker.New()
	sectionNormalizer := normalizer.New()
	sectionPartitioner := partitioner.New()
	expressionAnalyzer := analyzer.New()
	sectionExtractor := extractor.New(sectionTaker, sectionNormalizer, sectionPartitioner, expressionAnalyzer)
	variablesReplacer := replacer.New()
	funcCaller := caller.New()
	theParser := parser.New(doFileReader, commentCleaner, sectionExtractor, variablesReplacer, funcCaller)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}

			doFile, err := theParser.ParseFromFilename(tc.path)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if doFile != nil && tc.expected == nil {
				t.Errorf("expected nil, got %v", doFile)
			} else if doFile == nil && tc.expected != nil {
				t.Errorf("expected %v, got nil", tc.expected)
			} else if doFile != nil && tc.expected != nil {
				if doFile.Let.Variables != nil && tc.expected.Let.Variables == nil {
					for k, v := range doFile.Let.Variables {
						expectedVal, ok := tc.expected.Let.Variables[k]
						if !ok {
							t.Errorf("expected %v, got %v", expectedVal, v)
						}

						if !reflect.DeepEqual(expectedVal, v) {
							t.Errorf("expected %v, got %v", expectedVal, v)
							t.Errorf("expected %T, got %T", expectedVal, v)
						}
					}
				}

				if doFile.Do.Method != tc.expected.Do.Method {
					t.Errorf("expected %v, got %v", tc.expected.Do.Method, doFile.Do.Method)
				}

				if doFile.Do.URL != tc.expected.Do.URL {
					t.Errorf("expected %v, got %v", tc.expected.Do.URL, doFile.Do.URL)
				}

				if doFile.Do.Params != nil && tc.expected.Do.Params != nil {
					for k, v := range doFile.Do.Params {
						expectedVal, ok := tc.expected.Do.Params[k]
						if !ok {
							t.Errorf("expected %v, got %v", expectedVal, v)
						}

						if !reflect.DeepEqual(expectedVal, v) {
							t.Errorf("expected %v, got %v", expectedVal, v)
							t.Errorf("expected %T, got %T", expectedVal, v)
						}
					}
				}

				if doFile.Do.Headers != nil && tc.expected.Do.Headers != nil {
					for k, v := range doFile.Do.Headers {
						expectedVal, ok := tc.expected.Do.Headers[k]
						if !ok {
							t.Errorf("expected %v, got %v", expectedVal, v)
						}

						if !reflect.DeepEqual(expectedVal, v) {
							t.Errorf("expected %v, got %v", expectedVal, v)
							t.Errorf("expected %T, got %T", expectedVal, v)
						}
					}
				}

				if doFile.Do.Body != nil && tc.expected.Do.Body != nil {
					if *doFile.Do.Body != *tc.expected.Do.Body {
						t.Errorf("expected %v, got %v", *tc.expected.Do.Body, *doFile.Do.Body)
					}
				}
			}
		})
	}
}
