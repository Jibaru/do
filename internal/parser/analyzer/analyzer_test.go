package analyzer_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/types"
)

func TestAnalyzer_Analyze(t *testing.T) {
	testCases := []struct {
		name          string
		expressions   types.SectionExpressions
		expected      map[string]interface{}
		expectedError error
	}{
		{
			name: "success",
			expressions: types.SectionExpressions{
				"var1=1",
				"var2=\"hello\"",
				"var3=true",
				"var4=20.3",
				"var5=-12",
				"var6={\"key1\": 1, \"key2\": \"hello\"}",
				"var7=`something here`",
				"var8=\"=string=with=another=\"",
			},
			expected: map[string]interface{}{
				"var1": 1,
				"var2": "hello",
				"var3": true,
				"var4": 20.3,
				"var5": -12,
				"var6": map[string]interface{}{
					"key1": float64(1),
					"key2": "hello",
				},
				"var7": "something here",
				"var8": "=string=with=another=",
			},
		},
		{
			name: "error reading expression",
			expressions: types.SectionExpressions{
				"no equals",
			},
			expectedError: errors.New("error reading expression: [no equals]"),
		},
	}

	theAnalyzer := analyzer.New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := theAnalyzer.Analyze(tc.expressions)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			for key, value := range tc.expected {
				if !reflect.DeepEqual(value, result[key]) {
					t.Errorf("expected %v, got %v", value, result[key])
					t.Errorf("expected %T, got %T", value, result[key])
				}
			}
		})
	}
}
