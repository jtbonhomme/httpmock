package httpmock_test

import (
	"errors"
	"testing"

	"github.com/izysaas/go-kit/httpmock"
)

const ExampleText string = "The quick brown fox jumps over the lazy dog"

func TestTest(t *testing.T) {
	t.Run("assert", func(t *testing.T) {
		condition := true
		httpmock.Assert(t, condition, "expected condition to be true")
	})

	t.Run("ok", func(t *testing.T) {
		var condition error
		httpmock.Assert(t, condition == nil, "expected condition to be true")
		httpmock.OK(t, condition)
	})

	t.Run("not-nil", func(t *testing.T) {
		condition := errors.New("some error here")
		httpmock.NotNil(t, condition)
	})

	t.Run("nil", func(t *testing.T) {
		var condition error
		httpmock.Nil(t, condition)
	})

	t.Run("equals", func(t *testing.T) {
		tcs := []struct {
			message  string
			expected interface{}
			result   interface{}
		}{
			{
				message: "when expected is zero value",
			},
			{
				message:  "when expected is nil",
				expected: nil,
			},
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{"testing"},
				result:   struct{ test string }{"testing"},
			},
			{
				message:  "when expected and result are strings",
				expected: "testing",
				result:   "testing",
			},
		}
		for _, tc := range tcs {
			t.Log(tc.message)
			httpmock.Equals(t, tc.expected, tc.result)
		}
	})

	t.Run("not-zero", func(t *testing.T) {
		tcs := []struct {
			message  string
			expected interface{}
		}{
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{"testing"},
			},
			{
				message:  "when expected and result are strings",
				expected: "testing",
			},
			{
				message:  "when expected and result are integers",
				expected: 1,
			},
		}
		for _, tc := range tcs {
			t.Log(tc.message)
			httpmock.NotZero(t, tc.expected)
		}
	})

	t.Run("zero", func(t *testing.T) {
		tcs := []struct {
			message  string
			expected interface{}
		}{
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{},
			},
			{
				message:  "when expected and result are strings",
				expected: "",
			},
			{
				message:  "when expected and result are integers",
				expected: 0,
			},
		}
		for _, tc := range tcs {
			t.Log(tc.message)
			httpmock.Zero(t, tc.expected)
		}
	})

	t.Run("includes", func(t *testing.T) {
		result := ExampleText
		expected := "jumps"
		httpmock.Includes(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		httpmock.Includes(t, expected, resultList...)
	})

	t.Run("includes-i", func(t *testing.T) {
		result := ExampleText
		expected := "JUMPS"
		httpmock.IncludesI(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		httpmock.IncludesI(t, expected, resultList...)
	})

	t.Run("not-includes", func(t *testing.T) {
		result := ExampleText
		expected := "hippo"
		httpmock.NotIncludes(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		httpmock.NotIncludes(t, expected, resultList...)
	})

	t.Run("includes-slice", func(t *testing.T) {
		expected := []string{"B"}
		original := []string{"A", "B", "C"}
		httpmock.IncludesSlice(t, expected, original)

		expectedI := []int{5}
		originalI := []int{1, 2, 3, 4, 5, 6, 7}
		httpmock.IncludesSlice(t, expectedI, originalI)

		expectedE := []interface{}{5, "B"}
		originalE := []interface{}{1, 2, 3, 4, 5, 6, 7, "A", "B", "C"}
		httpmock.IncludesSlice(t, expectedE, originalE)
	})

	t.Run("includes-map", func(t *testing.T) {
		expected := map[string]string{"B": "B"}
		original := map[string]string{
			"A": "A",
			"B": "B",
			"C": "C",
		}
		httpmock.IncludesMap(t, expected, original)
	})
}
