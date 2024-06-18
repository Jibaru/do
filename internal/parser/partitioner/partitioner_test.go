package partitioner_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/types"
)

func TestPartitioner_Split(t *testing.T) {
	testCases := []struct {
		name     string
		content  types.NormalizedSectionContent
		expected types.SectionExpressions
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
	}

	thePartitioner := partitioner.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := thePartitioner.Split(tc.content)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
