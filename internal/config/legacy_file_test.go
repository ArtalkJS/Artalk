package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoadLegacyConfigPreservesSettings(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "artalk.yml")
	const old = `# 中文配置
app_key: fixture-key
allow_origins: ["https://blog.example.org"]
captcha:
  action_timeout: 42
moderator:
  pending_default: true
auth:
  enabled: true
notify:
  noise_mode: false
email:
  mail_subject_to_admin: "A comment needs moderation"
ip_region:
  enabled: true
`
	require.NoError(t, os.WriteFile(filename, []byte(old), 0600))
	t.Chdir(t.TempDir())

	parsed, err := config.NewFromFile(filename)
	require.NoError(t, err)
	require.Equal(t, "zh-CN", parsed.Locale)
	require.Equal(t, []string{"https://blog.example.org"}, parsed.TrustedDomains)
	require.Equal(t, 42, parsed.Captcha.ActionReset)
	require.Equal(t, config.TypeImage, parsed.Captcha.CaptchaType)
	require.True(t, parsed.AdminNotify.NoiseMode)
	require.Equal(t, "A comment needs moderation", parsed.AdminNotify.Email.MailSubject)
	require.Equal(t, "http://localhost:23366/api/v2/auth/:provider/callback", parsed.Auth.Callback)
	require.False(t, parsed.IPRegion.Enabled, "missing IP database disables the lookup")
}
