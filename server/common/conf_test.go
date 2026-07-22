package common

import (
	"testing"

	artalktest "github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlePluginURLs(t *testing.T) {
	app, err := artalktest.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, app.Cleanup())
	})

	app.Conf().TrustedDomains = []string{"https://trusted.example"}

	urls := handlePluginURLs(app.App, []string{
		"",
		"dist/plugins/local.js",
		"dist/plugins/local.js",
		"https://trusted.example/plugin.js",
		"https://untrusted.example/plugin.js",
	})

	assert.Equal(t, []string{
		"dist/plugins/local.js",
		"https://trusted.example/plugin.js",
	}, urls)
}
