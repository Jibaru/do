package replacer_test

import (
	"testing"

	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/types"
)

func TestVariablesReplacer_Replace(t *testing.T) {
	testCases := []struct {
		name         string
		doVariables  map[string]interface{}
		letVariables map[string]interface{}
		expected     map[string]interface{}
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
				},
				"body": types.String(`{"extra": $extra}`),
			},
			letVariables: map[string]interface{}{
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
			letVariables: map[string]interface{}{
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
			letVariables: map[string]interface{}{
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
	}

	varReplacer := replacer.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			varReplacer.Replace(tc.doVariables, tc.letVariables)

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
