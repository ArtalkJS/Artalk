package auth

import (
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCallbackOrigin(t *testing.T) {
	for _, tc := range []struct{ callback, origin string }{
		{"https://auth.example.com/behind/proxy", "https://auth.example.com"},
		{"http://localhost:23480", "http://localhost:23480"},
		{"", ""}, {"https://%", ""}, {"javascript:alert(1)", ""},
	} {
		t.Run(tc.callback, func(t *testing.T) {
			conf := &config.Config{}
			conf.Auth.Callback = tc.callback
			require.Equal(t, tc.origin, CallbackOrigin(conf))
		})
	}
}
