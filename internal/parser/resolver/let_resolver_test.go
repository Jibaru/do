package resolver_test

import (
	"reflect"
	"testing"

	"github.com/jibaru/do/internal/parser/resolver"
	"github.com/jibaru/do/internal/types"
)

func TestLetResolver_Resolve(t *testing.T) {
	testCases := []struct {
		name          string
		variables     *types.Sentences
		expected      *types.Sentences
		expectedError error
	}{
		{
			name: "success",
			variables: types.NewSentencesFromSlice([]types.Sentence{
				{
					Key:   "var1",
					Value: types.String("SOME_ARG"),
				},
				{
					Key: "var2",
					Value: types.Func{
						Name: "env",
						Args: []interface{}{
							types.ReferenceToVariable{Value: "var1"},
							types.String("default"),
						},
					},
				},
				{
					Key:   "var3",
					Value: types.ReferenceToVariable{Value: "var2"},
				},
			}),
			expected: types.NewSentencesFromSlice([]types.Sentence{
				{
					Key:   "var1",
					Value: types.String("SOME_ARG"),
				},
				{
					Key:   "var2",
					Value: types.String("default"),
				},
				{
					Key:   "var3",
					Value: types.String("default"),
				},
			}),
		},
		{
			name: "success with nested function",
			variables: types.NewSentencesFromSlice([]types.Sentence{
				{
					Key: "var1",
					Value: types.Func{
						Name: "env",
						Args: []interface{}{
							types.String("SOME_ARG"),
							types.String("default"),
						},
					},
				},
				{
					Key: "var2",
					Value: types.Func{
						Name: "env",
						Args: []interface{}{
							types.ReferenceToVariable{Value: "var1"},
							types.String("default"),
						},
					},
				},
				{
					Key:   "var3",
					Value: types.ReferenceToVariable{Value: "var2"},
				},
			}),
			expected: types.NewSentencesFromSlice([]types.Sentence{
				{
					Key:   "var1",
					Value: types.String("default"),
				},
				{
					Key:   "var2",
					Value: types.String("default"),
				},
				{
					Key:   "var3",
					Value: types.String("default"),
				},
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := resolver.NewLetResolver()
			resolvedVariables, err := r.Resolve(tc.variables)

			if err != nil && tc.expectedError == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("expected error %v, got no error", tc.expectedError)
			} else if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(resolvedVariables, tc.expected) {
				t.Errorf("expected: %v, got: %v", tc.expected, resolvedVariables)
			}
		})
	}
}
