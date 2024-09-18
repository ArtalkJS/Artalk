package config

import (
	"reflect"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_GetHashFuncByFrontendConf(t *testing.T) {
	tests := map[string]struct {
		config   *Config
		expected any
	}{
		"nil frontend conf": {
			config:   &Config{Frontend: nil},
			expected: utils.GetMD5Hash,
		},
		"empty frontend conf": {
			config:   &Config{Frontend: map[string]interface{}{}},
			expected: utils.GetMD5Hash,
		},
		"frontend conf without hash func": {
			config: &Config{Frontend: map[string]interface{}{
				"gravatar": false,
			}},
			expected: utils.GetMD5Hash,
		},
		"frontend conf with gravatar params but without hash func": {
			config: &Config{Frontend: map[string]interface{}{
				"gravatar": map[string]interface{}{"params": "a=1&b=2"},
			}},
			expected: utils.GetMD5Hash,
		},
		"frontend conf with gravatar params containing hash func": {
			config: &Config{Frontend: map[string]interface{}{
				"gravatar": map[string]interface{}{"params": "a=1&sha256=1&b=2"},
			}},
			expected: utils.GetSha256Hash,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expect := reflect.ValueOf(tc.expected).Pointer()
			actual := reflect.ValueOf(GetHashFuncByFrontendConf(tc.config)).Pointer()
			assert.Equal(t, expect, actual)
		})
	}
}
