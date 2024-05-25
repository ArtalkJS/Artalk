package artransfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isJsonArray(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"[1, 2, 3]", true},
		{"{\"a\": 1}", false},
		{"[]", true},
		{"", false},
		{"  [true]  ", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isJsonArray(test.input), "isJsonArray(%q)", test.input)
	}
}

func Test_convertJSONValuesToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`[{"a": 233}, {"b": true}, {"c": "233"}]`, `[{"a":"233"},{"b":"true"},{"c":"233"}]`},
		{`[{"a": 1.5}, {"b": false}]`, `[{"a":"1.5"},{"b":"false"}]`},
		{`[]`, `[]`},
	}

	for _, test := range tests {
		assert.JSONEq(t, test.expected, convertJSONValuesToString(test.input), "convertJSONValuesToString(%q)", test.input)
	}

	t.Run("Various type zero value", func(t *testing.T) {
		jsonStr := `[{"n":0},{"b":false},{"e":""},{"n":null},{"f1":0.0},{"f2":0.00}]`
		expected := `[{"n":"0"},{"b":"false"},{"e":""},{"n":""},{"f1":"0"},{"f2":"0"}]`

		result := convertJSONValuesToString(jsonStr)
		assert.Equal(t, expected, result)
	})
}

func Test_jsonDecodeFAS(t *testing.T) {
	tests := []struct {
		input       string
		expected    interface{}
		expectedErr bool
	}{
		{
			input: `[{"a": 233}, {"b": true}, {"c": "233"}]`,
			expected: &[]map[string]string{
				{"a": "233"},
				{"b": "true"},
				{"c": "233"},
			},
			expectedErr: false,
		},
		{
			input:       `{"a": 1}`, // not an array
			expected:    &[]map[string]string{},
			expectedErr: true,
		},
		{
			input:       `[{"a": "valid"}, {"b": "json"}]`,
			expected:    &[]map[string]string{{"a": "valid"}, {"b": "json"}},
			expectedErr: false,
		},
	}

	for _, test := range tests {
		var result []map[string]string
		err := jsonDecodeFAS(test.input, &result)
		if test.expectedErr {
			assert.Error(t, err, "jsonDecodeFAS(%q) expected error", test.input)
		} else {
			assert.NoError(t, err, "jsonDecodeFAS(%q) unexpected error", test.input)
			assert.Equal(t, test.expected, &result, "jsonDecodeFAS(%q)", test.input)
		}
	}
}
