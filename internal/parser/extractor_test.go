package parser_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/parser"
)

func TestSectionExtractor_Extract(t *testing.T) {
	testCases := []struct {
		name          string
		section       parser.Section
		rawContent    string
		expected      map[string]interface{}
		expectedError error
		normalizerFn  func(content string) (string, error)
	}{
		{
			name:       "success do section",
			section:    parser.DoSection,
			rawContent: `let{}do{method="GET";url="https://localhost:8080/api/v1/tests";}`,
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "https://localhost:8080/api/v1/tests",
			},
			normalizerFn: func(content string) (string, error) {
				return `method="GET";url="https://localhost:8080/api/v1/tests";`, nil
			},
		},
		{
			name:       "success let section",
			section:    parser.LetSection,
			rawContent: `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests"}`,
			expected: map[string]interface{}{
				"var1": 12,
				"var2": "text",
				"var3": false,
				"var4": 12.33,
			},
			normalizerFn: func(content string) (string, error) {
				return `var1=12;var2="text";var3=false;var4=12.33;`, nil
			},
		},
		{
			name:          "error no do block",
			section:       parser.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}`,
			expectedError: errors.New("no 'do' block found"),
			normalizerFn: func(content string) (string, error) {
				return "", nil
			},
		},
		{
			name:          "error missing opening brace after do",
			section:       parser.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do`,
			expectedError: errors.New("missing opening brace after 'do'"),
			normalizerFn: func(content string) (string, error) {
				return "", nil
			},
		},
		{
			name:          "error missing closing brace",
			section:       parser.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests"`,
			expectedError: errors.New("missing closing brace for 'do' block"),
			normalizerFn: func(content string) (string, error) {
				return "", nil
			},
		},
		{
			name:       "success let in string",
			section:    parser.LetSection,
			rawContent: `let{var1=12;var2="let";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expected: map[string]interface{}{
				"var1": 12,
				"var2": "let",
				"var3": false,
				"var4": 12.33,
			},
			normalizerFn: func(content string) (string, error) {
				return `var1=12;var2="let";var3=false;var4=12.33;`, nil
			},
		},
		{
			name:       "success do in string",
			section:    parser.DoSection,
			rawContent: `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://dolocalhost:8080/api/v1/tests";}`,
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "https://dolocalhost:8080/api/v1/tests",
			},
			normalizerFn: func(content string) (string, error) {
				return `method="GET";url="https://dolocalhost:8080/api/v1/tests";`, nil
			},
		},
		{
			name:          "error parsing JSON",
			section:       parser.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expectedError: errors.New("error parsing JSON value"),
			normalizerFn: func(content string) (string, error) {
				return `method="GET";url="https://localhost:8080/api/v1/tests";body={};`, nil
			},
		},
		{
			name:          "error parsing boolean value",
			section:       parser.DoSection,
			rawContent:    `let{var1=12;var2="text";var3=false;var4=12.33;}do{method="GET";url="https://localhost:8080/api/v1/tests";body={};}`,
			expectedError: errors.New("error parsing boolean value"),
			normalizerFn: func(content string) (string, error) {
				return `method="GET";url="https://localhost:8080/api/v1/tests";body={};`, nil
			},
		},
	}

	mockNormalizer := &parser.MockNormalizer{}
	sectionExtractor := parser.NewSectionExtractor(mockNormalizer)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockNormalizer.NormalizeFn = tc.normalizerFn

			result, err := sectionExtractor.Extract(tc.section, tc.rawContent)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			for key, value := range tc.expected {
				if result[key] != value {
					t.Errorf("expected %v, got %v", value, result[key])
					t.Errorf("expected %T, got %T", value, result[key])
				}
			}
		})
	}
}
