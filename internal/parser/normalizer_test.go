package parser_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/parser"
)

func TestDoSectionNormalizer_Normalize(t *testing.T) {
	testCases := []struct {
		name          string
		content       string
		expected      string
		expectedError error
	}{
		{
			name:     "success do section",
			content:  "    method = \"GET\";\n    url = \"http://    localhost:8080/api/todos/:id\";\n    params = {\"id\": \"$id\"};\n    query = {\n        \"id\": \"$id\"\n    };\n    headers = {\n        \"Authorization\": \"$token\",\n        \"Content-Type\": \"application/json\"\n};\n    body = `{\n\"extra\": \"something\"}\n`;",
			expected: "method=\"GET\";url=\"http://    localhost:8080/api/todos/:id\";params={\"id\":\"$id\"};query={\"id\":\"$id\"};headers={\"Authorization\":\"$token\",\"Content-Type\":\"application/json\"};body=`{\"extra\":\"something\"}`;",
		},
		{
			name:     "success let section",
			content:  "    var1 = 12; \n  var2 = \"text\"; \t   var3 = false;    var4 = 12.33;",
			expected: "var1=12;var2=\"text\";var3=false;var4=12.33;",
		},
		{
			name:          "error by empty content",
			content:       "",
			expectedError: errors.New("empty content"),
		},
		{
			name:          "error by empty spaced content",
			content:       "    ",
			expectedError: errors.New("empty content"),
		},
		{
			name:     "success let section with multiple ; and {}",
			content:  "    var1 = 12; \n  var2 = \"tex;;;t\"; \t   var3 = false;  var4 = 12.33;    var5=\"{;\";",
			expected: "var1=12;var2=\"tex;;;t\";var3=false;var4=12.33;var5=\"{;\";",
		},
	}

	normalizer := parser.NewNormalizer()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedContent, err := normalizer.Normalize(tc.content)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if normalizedContent != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, normalizedContent)
			}
		})
	}
}
