package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageTitleAutoFetchConfiguration(t *testing.T) {
	// Isolate both supported environment spellings from the developer's settings.
	for _, name := range []string{"ATK_PAGE_TITLE_AUTO_FETCH", "ATK_PAGE__TITLE__AUTO__FETCH"} {
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
	for _, tc := range []struct {
		name, yaml string
		want       bool
	}{
		{"omitted", "", false},
		{"enabled", "page_title_auto_fetch: true\n", true},
		{"disabled", "page_title_auto_fetch: false\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "artalk.yml")
			require.NoError(t, os.WriteFile(file, []byte("app_key: config-test-key\n"+tc.yaml), 0600))
			conf, err := NewFromFile(file)
			require.NoError(t, err)
			require.Equal(t, tc.want, conf.PageTitleAutoFetch)
			t.Setenv("ATK_PAGE_TITLE_AUTO_FETCH", "true")
			conf, err = NewFromFile(file)
			require.NoError(t, err)
			require.True(t, conf.PageTitleAutoFetch)
			t.Setenv("ATK_PAGE_TITLE_AUTO_FETCH", "false")
			conf, err = NewFromFile(file)
			require.NoError(t, err)
			require.False(t, conf.PageTitleAutoFetch)
		})
	}
}
