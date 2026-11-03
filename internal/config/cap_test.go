package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapConfiguration(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "artalk.yml")
	require.NoError(t, os.WriteFile(filename, []byte("locale: en\n"), 0600))
	conf, err := NewFromFile(filename)
	require.NoError(t, err)
	require.Equal(t, TypeImage, conf.Captcha.CaptchaType)
	t.Setenv("ATK_CAPTCHA_CAPTCHA_TYPE", "capjs")
	t.Setenv("ATK_CAPTCHA_CAPJS_KEY_ID", "site-id")
	t.Setenv("ATK_CAPTCHA_CAPJS_SECRET_KEY", "test-secret")
	t.Setenv("ATK_CAPTCHA_CAPJS_API_ENDPOINT", "https://cap.example/subpath/")
	conf, err = NewFromFile(filename)
	require.NoError(t, err)
	require.Equal(t, TypeCapJS, conf.Captcha.CaptchaType)
	require.Equal(t, CapJSConf{KeyID: "site-id", SecretKey: "test-secret", APIEndpoint: "https://cap.example/subpath/"}, conf.Captcha.CapJS)
	for _, field := range []string{"key_id", "secret_key", "api_endpoint"} {
		found := false
		for _, option := range OptionsMetaCache {
			if option.Path == "captcha.capjs."+field {
				require.True(t, option.AllowsSet)
				found = true
				break
			}
		}
		require.True(t, found, "missing metadata for %s", field)
	}
}
