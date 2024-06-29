package caller_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/types"
)

func TestCaller_Call(t *testing.T) {
	testCases := []struct {
		name                 string
		letVariables         map[string]interface{}
		doVariables          map[string]interface{}
		expectedLetVariables map[string]interface{}
		expectedDoVariables  map[string]interface{}
		expectedError        error
	}{
		{
			name: "success",
			letVariables: map[string]interface{}{
				"var1": types.String("value1"),
				"var2": types.EnvFunc{
					Arg1: "NO_EXISTS",
					Arg2: "default1",
				},
				"var3": types.FileFunc{
					Path: "/path/to/file",
				},
			},
			doVariables: map[string]interface{}{
				"var3": types.String("value3"),
				"var4": types.EnvFunc{
					Arg1: "NO_EXISTS",
					Arg2: "default2",
				},
				"var5": types.FileFunc{
					Path: "/path/to/file",
				},
			},
			expectedLetVariables: map[string]interface{}{
				"var1": types.String("value1"),
				"var2": types.String("default1"),
				"var3": types.File{Path: "/path/to/file"},
			},
			expectedDoVariables: map[string]interface{}{
				"var3": types.String("value3"),
				"var4": types.String("default2"),
				"var5": types.File{Path: "/path/to/file"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := caller.New()
			err := c.Call(tc.letVariables, tc.doVariables)
			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(tc.letVariables, tc.expectedLetVariables) {
				t.Errorf("expected letVariables %v, got %v", tc.expectedLetVariables, tc.letVariables)
			}
			if !reflect.DeepEqual(tc.doVariables, tc.expectedDoVariables) {
				t.Errorf("expected doVariables %v, got %v", tc.expectedDoVariables, tc.doVariables)
			}
		})
	}
}
