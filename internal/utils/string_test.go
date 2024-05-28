package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddQueryToURL(t *testing.T) {
	baseURL := "https://example.com"
	queryMap := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	resultURL := AddQueryToURL(baseURL, queryMap)
	expectedURL := "https://example.com?param1=value1&param2=value2"

	assert.Equal(t, resultURL, expectedURL)
}

func TestContainsStr(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	existing := "banana"
	nonExisting := "grape"

	assert.True(t, ContainsStr(slice, existing), "Expected true to be present in the slice")
	assert.False(t, ContainsStr(slice, nonExisting), "Expected false to not be present in the slice")
}

func TestRemoveDuplicates(t *testing.T) {
	input := []string{"apple", "banana", "cherry", "apple", "banana"}
	expected := []string{"apple", "banana", "cherry"}

	result := RemoveDuplicates(input)

	assert.Equal(t, len(result), len(expected))

	for i, val := range result {
		assert.Equal(t, expected[i], val, "Expected '%s' at index %d, but got '%s'", expected[i], i, val)
	}
}

func TestSplitAndTrimSpace(t *testing.T) {
	input := "apple, banana, cherry"
	sep := ","
	expected := []string{"apple", "banana", "cherry"}

	result := SplitAndTrimSpace(input, sep)
	assert.Equal(t, len(expected), len(result))

	for i, v := range result {
		assert.Equal(t, expected[i], v, "Expected '%s' at index %d, but got '%s'", expected[i], i, v)
	}
}

func TestRemoveBlankStrings(t *testing.T) {
	input := []string{"apple", "", "banana", "  ", "cherry", "  "}
	expected := []string{"apple", "banana", "cherry"}

	result := RemoveBlankStrings(input)
	assert.Equal(t, len(expected), len(result))

	for i, v := range result {
		assert.Equal(t, expected[i], v, "Expected '%s' at index %d, but got '%s'", expected[i], i, v)
	}
}

func TestTruncateString(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	length := 10
	expected := "abcdefghij"

	result := TruncateString(input, length)
	assert.Equal(t, expected, result)
}

func TestToString(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected string
	}{
		{true, "true"},
		{42, "42"},
		{"hello", "hello"},
	}

	for _, test := range tests {
		result := ToString(test.value)
		assert.Equal(t, test.expected, result, "For %v, expected %s, but got %s", test.value, test.expected, result)
	}
}
