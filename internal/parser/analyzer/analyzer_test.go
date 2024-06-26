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
		expected      *types.Sentences
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
				"var6={\"key1\": 1, \"key2\": \"hello\", \"key3\": {\"a\": x}, \"key4\": true, \"key5\": 20.3, \"key6\": -12, \"key7\": x, \"key8\": `backticks`}",
				"var7=`something here`",
				"var8=\"=string=with=another=\"",
				"var9=x",
				"var10=_y",
				"var11=env(\"OS_VAR\", \"default\")",
				"var12=file(\"/path/to/file\")",
			},
			expected: types.NewSentencesFromSlice([]types.Sentence{
				{
					Key:   "var1",
					Value: types.Int(1),
				},
				{
					Key:   "var2",
					Value: types.String("hello"),
				},
				{
					Key:   "var3",
					Value: types.Bool(true),
				},
				{
					Key:   "var4",
					Value: types.Float(20.3),
				},
				{
					Key:   "var5",
					Value: types.Int(-12),
				},
				{
					Key: "var6",
					Value: types.Map{
						"key1": types.Int(1),
						"key2": types.String("hello"),
						"key3": types.Map{
							"a": types.ReferenceToVariable{Value: "x"},
						},
						"key4": types.Bool(true),
						"key5": types.Float(20.3),
						"key6": types.Int(-12),
						"key7": types.ReferenceToVariable{Value: "x"},
						"key8": types.String("backticks"),
					},
				},
				{
					Key:   "var7",
					Value: types.String("something here"),
				},
				{
					Key:   "var8",
					Value: types.String("=string=with=another="),
				},
				{
					Key:   "var9",
					Value: types.ReferenceToVariable{Value: "x"},
				},
				{
					Key:   "var10",
					Value: types.ReferenceToVariable{Value: "_y"},
				},
				{
					Key:   "var11",
					Value: types.Func{Name: "env", Args: []interface{}{types.String("OS_VAR"), types.String("default")}},
				},
				{
					Key:   "var12",
					Value: types.Func{Name: "file", Args: []interface{}{types.String("/path/to/file")}},
				},
			}),
		},
		/*{
			// TODO: make sure this test is passing
			name: "success map with call",
			expressions: types.SectionExpressions{
				"var12={\"key1\": env(\"OS_VAR\", \"default2\")}",
			},
			expected: map[string]interface{}{
				"var12": types.Map{
					"key1": types.EnvFunc{Arg1: "OS_VAR", Arg2: "default2"},
				},
			},
		},*/
		{
			name: "error reading expression",
			expressions: types.SectionExpressions{
				"no equals",
			},
			expectedError: errors.New("error reading expression: [no equals]"),
		},
		{
			name: "error reserved keyword",
			expressions: types.SectionExpressions{
				"let=1",
			},
			expectedError: errors.New("reserved keyword let"),
		},
		{
			name: "error repeated key",
			expressions: types.SectionExpressions{
				"var1=1",
				"var1=2",
			},
			expectedError: errors.New("repeated key var1"),
		},
		{
			name: "error invalid value",
			expressions: types.SectionExpressions{
				"var1=1.1.1",
			},
			expectedError: errors.New("invalid value 1.1.1"),
		},
		{
			name: "error to map invalid value",
			expressions: types.SectionExpressions{
				"var1={\"key1\":: 1}",
			},
			expectedError: errors.New("invalid value : 1"),
		},
		{
			name: "error to map invalid value",
			expressions: types.SectionExpressions{
				"var1={\"key\": 1.1.1}",
			},
			expectedError: errors.New("invalid value 1.1.1"),
		},
		{
			name: "error invalid map in map value",
			expressions: types.SectionExpressions{
				"var1={\"key\": {\"key\": 1.1.1}}",
			},
			expectedError: errors.New("invalid value 1.1.1"),
		},
		{
			name: "error invalid map format",
			expressions: types.SectionExpressions{
				"var1={\"key\": 1, \"key\"}",
			},
			expectedError: errors.New("invalid value  \"key\""),
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

			if tc.expected == nil {
				return
			}

			for _, entry := range tc.expected.Entries() {
				key := entry.Key
				value := entry.Value

				expected, _ := result.Get(key)

				if !reflect.DeepEqual(value, expected) {
					t.Errorf("expected %v, got %v", value, expected)
					t.Errorf("expected %T, got %T", value, expected)
				}
			}
		})
	}
}
