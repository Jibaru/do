package cleaner_test

import (
	"testing"

	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/types"
)

func TestCleaner_Clean(t *testing.T) {
	testCases := []struct {
		name          string
		rawContent    types.FileReaderContent
		expected      types.CleanedContent
		expectedError error
	}{
		{
			name:       "no comments",
			rawContent: types.FileReaderContent("let {a = 1;}"),
			expected:   types.CleanedContent("let {a = 1;}"),
		},
		{
			name:       "one line comment",
			rawContent: types.FileReaderContent("let {a = 1;} // this is a comment"),
			expected:   types.CleanedContent("let {a = 1;} "),
		},
		{
			name:       "one line comment with quotes",
			rawContent: types.FileReaderContent("let {a = 1;} // this is a \"comment\""),
			expected:   types.CleanedContent("let {a = 1;} "),
		},
		{
			name:       "comments inside string",
			rawContent: types.FileReaderContent(`let {a = "//no remove";}`),
			expected:   types.CleanedContent(`let {a = "//no remove";}`),
		},
		{
			name:       "comments inside string with backticks",
			rawContent: types.FileReaderContent("let {a = `//no remove`;}`"),
			expected:   types.CleanedContent("let {a = `//no remove`;}`"),
		},
		{
			name:       "comment in next line",
			rawContent: types.FileReaderContent("let {a = 12;}\n// this is a comment"),
			expected:   types.CleanedContent("let {a = 12;}\n"),
		},
		{
			name:       "multiple lines with comments",
			rawContent: "let {\n    //var1 = 1;\n    var2 = \"hello\";\n //   var3 = true;\n    var4 = false; // Comment // //\n    var5 = var1; // Comment//\n}\n// Comment\ndo {\n    //method = \"GET\";\n    url = \"http://example.com/:id\";//Comment\n    params = {\n        \"id\": \"$var1\" // Comment with many words\n    };\n    headers = {\n        \"Content-Type\": \"application/json\",\n        \"X-Message\": \"$var2\",\n//        \"X-Var5\": var5\n    };\n    body = `{\n\t\t\t\"var1\": $var1,\n\t\t\t\"var2\": \"$var2\",\n\t\t\t//\"var3\": $var3,\n\t\t\t\"var4\": $var4,\n\t\t\t\"var5\": $var5\n\t\t}`;\n}\n//end",
			expected:   "let {\n    \n    var2 = \"hello\";\n \n    var4 = false; \n    var5 = var1; \n}\n\ndo {\n    \n    url = \"http://example.com/:id\";\n    params = {\n        \"id\": \"$var1\" \n    };\n    headers = {\n        \"Content-Type\": \"application/json\",\n        \"X-Message\": \"$var2\",\n\n    };\n    body = `{\n\t\t\t\"var1\": $var1,\n\t\t\t\"var2\": \"$var2\",\n\t\t\t//\"var3\": $var3,\n\t\t\t\"var4\": $var4,\n\t\t\t\"var5\": $var5\n\t\t}`;\n}\n",
		},
	}

	c := cleaner.New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cleanedContent, err := c.Clean(tc.rawContent)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if cleanedContent != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, cleanedContent)
			}
		})
	}
}
