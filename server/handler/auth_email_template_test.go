package handler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestEmailVerificationTemplateFallback(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "verification.html")
	const content = "Your code: {{code}}"
	require.NoError(t, os.WriteFile(file, []byte(content), 0600))
	for _, tc := range []struct{ name, path, want string }{
		{"custom", file, content},
		{"default", "default", ""},
		{"missing", filepath.Join(dir, "missing.html"), ""},
		{"unreadable", dir, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, loadAuthEmailVerifyTemplate(&config.AuthEmailConf{VerifyTpl: tc.path}))
		})
	}
}
