package replacer_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/types"
)

func TestVariablesReplacer_Replace(t *testing.T) {
	testCases := []struct {
		name          string
		doVariables   map[string]interface{}
		letVariables  types.Map
		expected      map[string]interface{}
		expectedError error
	}{
		{
			name: "success",
			doVariables: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "$id"},
				"query":  types.Map{"ok": "$isOk"},
				"headers": types.Map{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
					"Extra":         types.NewReferenceToVariable("extra"),
				},
				"body": types.String(`{"extra": $extra}`),
			},
			letVariables: types.Map{
				"id":    types.String("123"),
				"isOk":  types.Bool(true),
				"token": types.String("random123"),
				"extra": types.Int(902),
				"other": types.Int(902),
			},
			expected: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "123"},
				"query":  types.Map{"ok": "true"},
				"headers": types.Map{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
					"Extra":         types.Int(902),
				},
				"body": types.String(`{"extra": 902}`),
			},
		},
		{
			name: "success no let variables",
			doVariables: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "$id"},
				"query":  types.Map{"ok": "$isOk"},
				"headers": types.Map{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
				},
				"body": types.String(`{"extra": $extra}`),
			},
			letVariables: nil,
			expected: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "$id"},
				"query":  types.Map{"ok": "$isOk"},
				"headers": types.Map{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
				},
				"body": types.String(`{"extra": $extra}`),
			},
		},
		{
			name:        "success no do variables",
			doVariables: nil,
			letVariables: types.Map{
				"id":    types.String("123"),
				"isOk":  types.Bool(true),
				"token": types.String("random123"),
			},
			expected: nil,
		},
		{
			name: "success no needs replacements",
			doVariables: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "123"},
				"query":  types.Map{"ok": "true"},
				"headers": types.Map{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
				},
				"body": types.String(`{"extra": 902}`),
			},
			letVariables: types.Map{
				"id":    types.String("123"),
				"isOk":  types.Bool(true),
				"token": types.String("random123"),
				"extra": types.Int(902),
			},
			expected: map[string]interface{}{
				"method": types.String("GET"),
				"url":    types.String("http://localhost:8080/api/todos/:id"),
				"params": types.Map{"id": "123"},
				"query":  types.Map{"ok": "true"},
				"headers": types.Map{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
				},
				"body": types.String(`{"extra": 902}`),
			},
		},
		{
			name:        "error invalid let variables",
			doVariables: map[string]interface{}{},
			letVariables: types.Map{
				"id": types.Map{"id": "123"},
				"func": types.EnvFunc{
					Arg1: "TEST_1",
					Arg2: "DEFAULT",
				},
				"reference": types.NewReferenceToVariable("id"),
			},
			expectedError: errors.New("let variables must have basic types values"),
		},
		{
			name: "error reference to variable not found in do section",
			letVariables: types.Map{
				"var1": types.String("123"),
			},
			doVariables: map[string]interface{}{
				"method": types.NewReferenceToVariable("var2"),
			},
			expectedError: errors.New("reference to variable for key method not found: var2"),
		},
		{
			name: "error reference to variable not found in do section nested map",
			letVariables: types.Map{
				"var1": types.String("123"),
			},
			doVariables: map[string]interface{}{
				"map_variable": types.Map{
					"method": types.NewReferenceToVariable("var2"),
				},
			},
			expectedError: errors.New("reference to variable for key method not found: var2"),
		},
	}

	varReplacer := replacer.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := varReplacer.Replace(tc.doVariables, tc.letVariables)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			for key, value := range tc.doVariables {
				switch val := value.(type) {
				case string:
					if val != tc.expected[key] {
						t.Errorf("expected %v, got %v", tc.expected[key], val)
						t.Errorf("expected %T, got %T", tc.expected[key], val)
					}
				case map[string]interface{}:
					for k, v := range val {
						if v != tc.expected[key].(map[string]interface{})[k] {
							t.Errorf("expected %v, got %v", tc.expected[key].(map[string]interface{})[k], v)
							t.Errorf("expected %T, got %T", tc.expected[key].(map[string]interface{})[k], v)
						}
					}
				}
			}
		})
	}

}
