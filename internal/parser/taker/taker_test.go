package taker_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/types"
)

func TestTaker_Take(t *testing.T) {
	testCases := []struct {
		name          string
		section       types.Section
		text          types.FileReaderContent
		expected      types.RawSectionContent
		expectedError error
	}{
		{
			name:     "success let section",
			section:  types.LetSection,
			text:     " let   {    var1 = 12; \n  var2 = \"text\"; \t   var3 = false;    var4 = 12.33;}",
			expected: "    var1 = 12; \n  var2 = \"text\"; \t   var3 = false;    var4 = 12.33;",
		},
		{
			name:     "success let section with multiple ; and {}",
			section:  types.LetSection,
			text:     "let{var1=12;var2=\"tex;;;t\";var3=false;var4=12.33;var5=\"{;\";}do{method=\"GET\";}",
			expected: "var1=12;var2=\"tex;;;t\";var3=false;var4=12.33;var5=\"{;\";",
		},
		{
			name:     "success do section",
			section:  types.DoSection,
			text:     "let {\n    var1 = 1;\n    var2 = \"hello\";\n    var3 = true;\n    var4 = false;\n}\n\ndo {\n    method = \"GET\";\n    url = \"http://example.com/:id\";\n    params = {\n        \"id\": \"$var1\"\n    };\n    headers = {\n        \"Content-Type\": \"application/json\",\n        \"X-Message\": \"$var2\"\n    };\n    body = `{\n        \"var1\": $var1,\n        \"var2\": \"$var2\",\n        \"var3\": $var3,\n        \"var4\": $var4\n    }`;\n}",
			expected: "\n    method = \"GET\";\n    url = \"http://example.com/:id\";\n    params = {\n        \"id\": \"$var1\"\n    };\n    headers = {\n        \"Content-Type\": \"application/json\",\n        \"X-Message\": \"$var2\"\n    };\n    body = `{\n        \"var1\": $var1,\n        \"var2\": \"$var2\",\n        \"var3\": $var3,\n        \"var4\": $var4\n    }`;\n",
		},
		{
			name:          "no do block error",
			section:       types.DoSection,
			text:          "let {var1=23}",
			expectedError: errors.New("no block found"),
		},
		{
			name:          "no let block error",
			section:       types.LetSection,
			text:          "do {method=\"GET\"}",
			expectedError: errors.New("no block found"),
		},
		{
			name:          "no let block error empty content",
			section:       types.LetSection,
			text:          "let {}",
			expectedError: errors.New("no block found"),
		},
		{
			name:          "missing opening brace error",
			section:       types.LetSection,
			text:          "let var1=12; var2=\"text\"; var3=false; var4=12.33;}",
			expectedError: errors.New("missing opening brace"),
		},
		{
			name:          "missing closing brace error",
			section:       types.LetSection,
			text:          "let {var1=12; var2=\"text\"; var3=false; var4=12.33;",
			expectedError: errors.New("missing closing brace"),
		},
	}

	theTaker := taker.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := theTaker.Take(tc.section, tc.text)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if content != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, content)
			}
		})
	}
}
