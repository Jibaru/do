package extractor_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/types"
)

func TestExtractor_Extract(t *testing.T) {
	testCases := []struct {
		name          string
		section       types.Section
		rawContent    types.FileReaderContent
		expected      map[string]interface{}
		expectedError error
		takeFn        func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error)
		normalizerFn  func(content types.RawSectionContent) (types.NormalizedSectionContent, error)
		splitFn       func(content types.NormalizedSectionContent) (types.SectionExpressions, error)
		analyzeFn     func(expressions types.SectionExpressions) (map[string]interface{}, error)
	}{
		{
			name:       "success do section",
			section:    types.DoSection,
			rawContent: "let{}do{method=\"GET\";url=\"https://localhost:8080/api/v1/tests\";params={\"id\":12};headers={\"Authorization\":\"Bearer token\"};body=`{\"extra\":true}`;}",
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "https://localhost:8080/api/v1/tests",
				"params": map[string]interface{}{
					"id": float64(12),
				},
				"headers": map[string]interface{}{
					"Authorization": "Bearer token",
				},
				"body": `{"extra":true}`,
			},
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "method=\"GET\";url=\"https://localhost:8080/api/v1/tests\";params={\"id\":12};headers={\"Authorization\":\"Bearer token\"};body=`{\"extra\":true}`;", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return types.NormalizedSectionContent(content), nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"method=\"GET\"",
					"url=\"https://localhost:8080/api/v1/tests\"",
					"params={\"id\":12}",
					"headers={\"Authorization\":\"Bearer token\"}",
					"body=`{\"extra\":true}`",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return map[string]interface{}{
					"method": "GET",
					"url":    "https://localhost:8080/api/v1/tests",
					"params": map[string]interface{}{
						"id": float64(12),
					},
					"headers": map[string]interface{}{
						"Authorization": "Bearer token",
					},
					"body": `{"extra":true}`,
				}, nil
			},
		},
		{
			name:       "success let section",
			section:    types.LetSection,
			rawContent: `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests"}`,
			expected: map[string]interface{}{
				"var1": 12,
				"var2": "text",
				"var3": false,
				"var4": 12.33,
			},
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "var1=12;var2=\"text\";var3=false;var4=12.33;", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `var1=12;var2="text";var3=false;var4=12.33;`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"var1=12",
					"var2=\"text\"",
					"var3=false",
					"var4=12.33",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return map[string]interface{}{
					"var1": 12,
					"var2": "text",
					"var3": false,
					"var4": 12.33,
				}, nil
			},
		},
		{
			name:          "error no do block",
			section:       types.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}`,
			expectedError: errors.New("no block found"),
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return "", nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return nil, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return nil, nil
			},
		},
		{
			name:          "error missing opening brace after do",
			section:       types.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do`,
			expectedError: errors.New("missing opening brace"),
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return "", nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return nil, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return nil, nil
			},
		},
		{
			name:          "error missing closing brace",
			section:       types.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests"`,
			expectedError: errors.New("missing closing brace"),
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return "", nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return nil, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return nil, nil
			},
		},
		{
			name:       "success let in string",
			section:    types.LetSection,
			rawContent: `let{var1=12;var2="let";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expected: map[string]interface{}{
				"var1": 12,
				"var2": "let",
				"var3": false,
				"var4": 12.33,
			},
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "var1=12;var2=\"let\";var3=false;var4=12.33;", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `var1=12;var2="let";var3=false;var4=12.33;`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"var1=12",
					"var2=\"let\"",
					"var3=false",
					"var4=12.33",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return map[string]interface{}{
					"var1": 12,
					"var2": "let",
					"var3": false,
					"var4": 12.33,
				}, nil
			},
		},
		{
			name:       "success do in string",
			section:    types.DoSection,
			rawContent: `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://dolocalhost:8080/api/v1/tests";}`,
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "https://dolocalhost:8080/api/v1/tests",
			},
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "method=\"GET\";url=\"https://dolocalhost:8080/api/v1/tests\";", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `method="GET";url="https://dolocalhost:8080/api/v1/tests";`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"method=\"GET\"",
					"url=\"https://dolocalhost:8080/api/v1/tests\"",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return map[string]interface{}{
					"method": "GET",
					"url":    "https://dolocalhost:8080/api/v1/tests",
				}, nil
			},
		},
		{
			name:          "error parsing JSON",
			section:       types.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expectedError: errors.New("error parsing JSON value"),
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "method=\"GET\";url=\"https://localhost:8080/api/v1/tests\";body={};", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `method="GET";url="https://localhost:8080/api/v1/tests";body={};`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"method=\"GET\"",
					"url=\"https://localhost:8080/api/v1/tests\"",
					"body={}",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return nil, nil
			},
		},
		{
			name:          "error parsing boolean value",
			section:       types.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expectedError: errors.New("error parsing boolean value"),
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "method=\"GET\";url=\"https://localhost:8080/api/v1/tests\";body={};", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `method="GET";url="https://localhost:8080/api/v1/tests";body={};`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"method=\"GET\"",
					"url=\"https://localhost:8080/api/v1/tests\"",
					"body={}",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return nil, nil
			},
		},
		{
			name:       "success let section with multiple ; and {}",
			section:    types.LetSection,
			rawContent: "let{var1=12;var2=\"tex;;;t\";var3=false;var4=12.33;var5=\"{;\";}do{method=\"GET\";}",
			expected: map[string]interface{}{
				"var1": 12,
				"var2": "tex;;;t",
				"var3": false,
				"var4": 12.33,
				"var5": "{;",
			},
			takeFn: func(section types.Section, text types.FileReaderContent) (types.RawSectionContent, error) {
				return "var1=12;var2=\"tex;;;t\";var3=false;var4=12.33;var5=\"{;\";", nil
			},
			normalizerFn: func(content types.RawSectionContent) (types.NormalizedSectionContent, error) {
				return `var1=12;var2="tex;;;t";var3=false;var4=12.33;var5="{;";`, nil
			},
			splitFn: func(content types.NormalizedSectionContent) (types.SectionExpressions, error) {
				return types.SectionExpressions{
					"var1=12",
					"var2=\"tex;;;t\"",
					"var3=false",
					"var4=12.33",
					"var5=\"{;\"",
				}, nil
			},
			analyzeFn: func(expressions types.SectionExpressions) (map[string]interface{}, error) {
				return map[string]interface{}{
					"var1": 12,
					"var2": "tex;;;t",
					"var3": false,
					"var4": 12.33,
					"var5": "{;",
				}, nil
			},
		},
	}

	normalizerMock := &normalizer.Mock{}
	partitionerMock := &partitioner.Mock{}
	analyzerMock := &analyzer.Mock{}
	takerMock := &taker.Mock{}
	sectionExtractor := extractor.New(takerMock, normalizerMock, partitionerMock, analyzerMock)

	isMap := func(i interface{}) bool {
		_, ok := i.(map[string]interface{})
		return ok
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			takerMock.TakeFn = tc.takeFn
			normalizerMock.NormalizeFn = tc.normalizerFn
			partitionerMock.SplitFn = tc.splitFn
			analyzerMock.AnalyzeFn = tc.analyzeFn

			result, err := sectionExtractor.Extract(tc.section, tc.rawContent)

			if err != nil && ((tc.expectedError == nil) || (err.Error() != tc.expectedError.Error())) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			for key, value := range tc.expected {
				if isMap(result[key]) && isMap(value) {
					for k, v := range value.(map[string]interface{}) {
						if result[key].(map[string]interface{})[k] != v {
							t.Errorf("expected %v, got %v", v, result[key].(map[string]interface{})[k])
							t.Errorf("expected %T, got %T", v, result[key].(map[string]interface{})[k])
						}
					}
					continue
				}

				if result[key] != value {
					t.Errorf("expected %v, got %v", value, result[key])
					t.Errorf("expected %T, got %T", value, result[key])
				}
			}
		})
	}
}
