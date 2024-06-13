package parser_test

import (
	"testing"

	"github.com/jibaru/do/internal/parser"
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
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "$id"},
				"query":  map[string]interface{}{"ok": "$isOk"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": $extra}`,
			},
			letVariables: map[string]interface{}{
				"id":    "123",
				"isOk":  true,
				"token": "random123",
				"extra": 902,
			},
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "123"},
				"query":  map[string]interface{}{"ok": "true"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": 902}`,
			},
		},
		{
			name: "success no let variables",
			doVariables: map[string]interface{}{
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "$id"},
				"query":  map[string]interface{}{"ok": "$isOk"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": $extra}`,
			},
			letVariables: nil,
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "$id"},
				"query":  map[string]interface{}{"ok": "$isOk"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer $token",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": $extra}`,
			},
		},
		{
			name:        "success no do variables",
			doVariables: nil,
			letVariables: map[string]interface{}{
				"id":    "123",
				"isOk":  true,
				"token": "random123",
			},
			expected: nil,
		},
		{
			name: "success no needs replacements",
			doVariables: map[string]interface{}{
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "123"},
				"query":  map[string]interface{}{"ok": "true"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": 902}`,
			},
			letVariables: map[string]interface{}{
				"id":    "123",
				"isOk":  true,
				"token": "random123",
				"extra": 902,
			},
			expected: map[string]interface{}{
				"method": "GET",
				"url":    "http://localhost:8080/api/todos/:id",
				"params": map[string]interface{}{"id": "123"},
				"query":  map[string]interface{}{"ok": "true"},
				"headers": map[string]interface{}{
					"Authorization": "Bearer random123",
					"Content-Type":  "application/json",
				},
				"body": `{"extra": 902}`,
			},
		},
	}

	replacer := parser.NewVariablesReplacer()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			replacer.Replace(tc.doVariables, tc.letVariables)

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
