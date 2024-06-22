package partitioner_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/types"
)

func TestPartitioner_Split(t *testing.T) {
	testCases := []struct {
		name          string
		content       types.NormalizedSectionContent
		expected      types.SectionExpressions
		expectedError error
	}{
		{
			name:    "success",
			content: "var1=12;var2=\"text\";var3=false;var4=12.33;var5=\";\";var6=`something`;var7={\"key1\": 1, \"key2\": \"hello\"};",
			expected: types.SectionExpressions{
				"var1=12",
				"var2=\"text\"",
				"var3=false",
				"var4=12.33",
				"var5=\";\"",
				"var6=`something`",
				"var7={\"key1\": 1, \"key2\": \"hello\"}",
			},
		},
		{
			name:          "error empty part",
			content:       "var1=2;;",
			expectedError: errors.New("empty part"),
		},
	}

	thePartitioner := partitioner.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := thePartitioner.Split(tc.content)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
