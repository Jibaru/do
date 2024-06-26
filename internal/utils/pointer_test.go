package utils

import "testing"

func TestPtr(t *testing.T) {
	integerVal := 1
	stringVal := "string"
	boolVal := true
	floatVal := 1.1
	mapVal := map[string]interface{}{"key": "value"}

	testCases := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "int",
			value:    integerVal,
			expected: &integerVal,
		},
		{
			name:     "string",
			value:    stringVal,
			expected: &stringVal,
		},
		{
			name:     "bool",
			value:    boolVal,
			expected: &boolVal,
		},
		{
			name:     "float",
			value:    floatVal,
			expected: &floatVal,
		},
		{
			name:     "map",
			value:    mapVal,
			expected: &mapVal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Ptr(tc.value)

			if *result != *tc.expected.(*int) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
