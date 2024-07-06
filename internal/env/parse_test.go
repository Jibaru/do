package env_test

import (
	"errors"
	"testing"

	"github.com/jibaru/do/internal/env"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		filepath      string
		expected      map[string]string
		expectedError error
	}{
		{
			name:     "success",
			filepath: "testdata/01_success.env",
			expected: map[string]string{
				"var1":  "content",
				"var2":  "true",
				"var3":  "42",
				"VAR_4": "3.14",
				"VAR_5": "content",
				"VAR_6": "content",
				"VAR_7": "=====content=====",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			envMap, err := env.Parse(tc.filepath)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if len(envMap) != len(tc.expected) {
				t.Errorf("expected %d env vars, got: %d", len(tc.expected), len(envMap))
			}

			for key, value := range tc.expected {
				if envMap[key] != value {
					t.Errorf("expected %s=%s, got: %s=%s", key, value, key, envMap[key])
				}
			}
		})
	}
}
