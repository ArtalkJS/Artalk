/**
 * List of sensitive config paths.
 *
 * (which should be hidden in the UI)
 */
export const SensitiveConfigPaths = [
  'app_key',
  'admin_notify.ding_talk.secret',
  'admin_notify.line.channel_access_token',
  'admin_notify.line.channel_secret',
  'admin_notify.lark.webhook_url',
  'admin_notify.slack.oauth_token',
  'admin_notify.telegram.api_token',
  'admin_notify.webhook.url',
  'admin.notify.bark.server',
  'auth.apple.client_secret',
  'auth.auth0.client_secret',
  'auth.discord.client_secret',
  'auth.facebook.client_secret',
  'auth.gitea.client_secret',
  'auth.github.client_secret',
  'auth.gitlab.client_secret',
  'auth.google.client_secret',
  'auth.line.client_secret',
  'auth.mastodon.client_secret',
  'auth.microsoft.client_secret',
  'auth.patreon.client_secret',
  'auth.slack.client_secret',
  'auth.tiktok.client_secret',
  'auth.twitter.client_secret',
  'auth.wechat.client_secret',
  'auth.steam.api_key',
  'captcha.geetest.captcha_key',
  'captcha.hcaptcha.secret_key',
  'captcha.recaptcha.secret_key',
  'captcha.turnstile.secret_key',
  'db.password',
  'email.ali_dm.access_key_secret',
  'email.smtp.password',
]

export function isSensitiveConfigPath(path: string) {
  return SensitiveConfigPaths.includes(path)
}
