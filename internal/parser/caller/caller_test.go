package caller_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

func TestCaller_Call(t *testing.T) {
	uuid := "80aaa8e2-e2b9-4bd5-8124-4003d4a528df"
	now := time.Now().UTC()

	testCases := []struct {
		name          string
		variables     map[string]interface{}
		expected      map[string]interface{}
		expectedError error
	}{
		{
			name: "success",
			variables: map[string]interface{}{
				"var1": types.String("value1"),
				"var2": types.Func{
					Name: "env",
					Args: []interface{}{
						types.String("NO_EXISTS"),
						types.String("default1"),
					},
				},
				"var3": types.Func{
					Name: "file",
					Args: []interface{}{
						types.String("/path/to/file"),
					},
				},
				"var4": types.Func{
					Name: "uuid",
				},
				"var5": types.Func{
					Name: "date",
					Args: []interface{}{
						types.String("ISO8601"),
					},
				},
			},
			expected: map[string]interface{}{
				"var1": types.String("value1"),
				"var2": types.String("default1"),
				"var3": types.File{Path: "/path/to/file"},
				"var4": types.String(uuid),
				"var5": types.String(now.Format(time.RFC3339)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := caller.New(
				utils.NewFixedUuidFactory(uuid),
				utils.NewFixedDateFactory(now),
			)
			err := c.Call(tc.variables)
			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(tc.variables, tc.expected) {
				t.Errorf("expected letVariables %v, got %v", tc.expected, tc.variables)
			}
		})
	}
}
