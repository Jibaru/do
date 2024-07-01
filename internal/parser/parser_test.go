package parser_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/resolver"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

func TestParser_FromFilename(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expected      *types.DoFile
		expectedError error
		FileReaderFn  func(filename string) (types.FileReaderContent, error)
		CleanerFn     func(content types.FileReaderContent) (types.CleanedContent, error)
		ExtractorFn   func(section types.Section, content types.CleanedContent) (*types.Sentences, error)
		ReplacerFn    func(doVariables map[string]interface{}, letVariables types.Map) error
		CallerFn      func(variables map[string]interface{}) error
		ResolverFn    func(variables *types.Sentences) (*types.Sentences, error)
	}{
		{
			name:     "success",
			filename: "valid.do",
			expected: &types.DoFile{
				Let: types.Let{
					Variables: nil,
				},
				Do: types.Do{
					Method: types.String("GET"),
					URL:    types.String("http://localhost:8080"),
				},
			},
			ExtractorFn: func(section types.Section, content types.CleanedContent) (*types.Sentences, error) {
				if section == types.DoSection {
					return types.NewSentencesFromSlice([]types.Sentence{
						{
							Key:   "method",
							Value: types.String("GET"),
						},
						{
							Key:   "url",
							Value: types.String("http://localhost:8080"),
						},
					}), nil
				}

				return nil, nil
			},
		},
	}

	fileReader := &reader.Mock{}
	commentCleaner := &cleaner.Mock{}
	sectionExtractor := &extractor.Mock{}
	varReplacer := &replacer.Mock{}
	funcCaller := &caller.Mock{}
	letResolver := &resolver.Mock{}

	p := parser.New(fileReader, commentCleaner, sectionExtractor, varReplacer, funcCaller, letResolver)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.FileReaderFn == nil {
				tc.FileReaderFn = func(filename string) (types.FileReaderContent, error) {
					return "", nil
				}
			}

			if tc.CleanerFn == nil {
				tc.CleanerFn = func(content types.FileReaderContent) (types.CleanedContent, error) {
					return "", nil
				}
			}

			if tc.ExtractorFn == nil {
				tc.ExtractorFn = func(section types.Section, content types.CleanedContent) (*types.Sentences, error) {
					return nil, nil
				}
			}

			if tc.ReplacerFn == nil {
				tc.ReplacerFn = func(doVariables map[string]interface{}, letVariables types.Map) error {
					return nil
				}
			}

			if tc.CallerFn == nil {
				tc.CallerFn = func(variables map[string]interface{}) error {
					return nil
				}
			}

			if tc.ResolverFn == nil {
				tc.ResolverFn = func(variables *types.Sentences) (*types.Sentences, error) {
					return nil, nil
				}
			}

			fileReader.ReadFn = tc.FileReaderFn
			commentCleaner.CleanFn = tc.CleanerFn
			sectionExtractor.ExtractFn = tc.ExtractorFn
			varReplacer.ReplaceFn = tc.ReplacerFn
			funcCaller.CallFn = tc.CallerFn
			letResolver.ResolveFn = tc.ResolverFn

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
