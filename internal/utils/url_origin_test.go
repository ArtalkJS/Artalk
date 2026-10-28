package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLOrigin(t *testing.T) {
	for _, tc := range []struct{ raw, origin string }{
		{" https://blog.test/path?x=1#anchor ", "https://blog.test"},
		{"http://localhost:8080/path", "http://localhost:8080"},
		{"http://[::1]:8080/path", "http://[::1]:8080"},
		{"https://blog.test@untrusted.test/path", "https://untrusted.test"},
	} {
		origin, err := URLOrigin(tc.raw)
		assert.NoError(t, err)
		assert.Equal(t, tc.origin, origin)
	}
	for _, raw := range []string{"", "*", "https://*.test", "//blog.test/path", "file:///tmp/page", "https://", "https://bad%2f.test"} {
		_, err := URLOrigin(raw)
		assert.Error(t, err, raw)
	}
}
