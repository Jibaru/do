package resolver_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/jibaru/do/internal/parser/resolver"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

func TestLetResolver_Resolve(t *testing.T) {
	uuid := "80aaa8e2-e2b9-4bd5-8124-4003d4a528df"
	now := time.Now().UTC()

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
				{
					Key: "var4",
					Value: types.Func{
						Name: "uuid",
					},
				},
				{
					Key: "var5",
					Value: types.Func{
						Name: "date",
						Args: []interface{}{
							types.String("ISO8601"),
						},
					},
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
				{
					Key:   "var4",
					Value: types.String(uuid),
				},
				{
					Key:   "var5",
					Value: types.String(now.Format(time.RFC3339)),
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
			r := resolver.NewLetResolver(
				utils.NewFixedUuidFactory(uuid),
				utils.NewFixedDateFactory(now),
			)
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
