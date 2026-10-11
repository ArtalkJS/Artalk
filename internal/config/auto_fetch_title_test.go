package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoFetchTitleConfiguration(t *testing.T) {
	// Isolate both supported environment spellings from the developer's settings.
	for _, name := range []string{"ATK_AUTO_FETCH_TITLE", "ATK_AUTO__FETCH__TITLE"} {
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
	for _, tc := range []struct {
		name, yaml string
		want       bool
	}{
		{"omitted", "", false},
		{"enabled", "auto_fetch_title: true\n", true},
		{"disabled", "auto_fetch_title: false\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "artalk.yml")
			require.NoError(t, os.WriteFile(file, []byte("app_key: config-test-key\n"+tc.yaml), 0600))
			conf, err := NewFromFile(file)
			require.NoError(t, err)
			require.Equal(t, tc.want, conf.AutoFetchTitle)
			t.Setenv("ATK_AUTO_FETCH_TITLE", "true")
			conf, err = NewFromFile(file)
			require.NoError(t, err)
			require.True(t, conf.AutoFetchTitle)
			t.Setenv("ATK_AUTO_FETCH_TITLE", "false")
			conf, err = NewFromFile(file)
			require.NoError(t, err)
			require.False(t, conf.AutoFetchTitle)
		})
	}
}
