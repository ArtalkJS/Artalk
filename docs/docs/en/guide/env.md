# Environment Variables

Artalk supports configuration changes via environment variables, simplifying containerized cloud deployments such as Docker.

Environment variable names start with `ATK_`, in all uppercase, corresponding to each node in the [configuration file](./backend/config.md).

Environment variables take precedence over configuration files. Verify if the configuration has taken effect using the following command:
  
```bash
artalk config
```

To deploy Artalk using Docker Compose, add environment variables in the `compose.yml` file, for instance:

```yaml
version: '3.8'
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    restart: unless-stopped
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
    environment:
      - TZ=America/New_York
      - ATK_LOCALE=en
      - ATK_SITE_DEFAULT=Artalk's Blog
      - ATK_SITE_URL=https://example.com
      - ATK_TRUSTED_DOMAINS=https://dev.example.com https://localhost:8080
      - ATK_ADMIN_USERS_0_NAME=admin
      - ATK_ADMIN_USERS_0_EMAIL=admin@example.org
      - ATK_ADMIN_USERS_0_PASSWORD=(bcrypt)$2y$10$ti4vZYIrxVN8rLcYXVgXCO.GJND0dyI49r7IoF3xqIx8bBRmIBZRm
      - ATK_ADMIN_USERS_0_BADGE_NAME=Administrator
      - ATK_ADMIN_USERS_0_BADGE_COLOR=#0083FF
```

When the variable is an array, set array values using space-separated strings or numerical indices, for example:

```
ATK_TRUSTED_DOMAINS="https://a.com https://b.com"
ATK_TRUSTED_DOMAINS_0="https://a.com"
```

<style>
.env-table table {
  table-layout: fixed;
  word-break: break-word;
  
  th:first-of-type, td:first-of-type {
    width: 30%;
  }
  th:nth-of-type(2), td:nth-of-type(2) {
    width: 20%;
  }
  th:nth-of-type(3), td:nth-of-type(3) {
    width: 20%;
  }
  th:nth-of-type(4), td:nth-of-type(4) {
    width: 30%;
  }
}
</style>

<div class="env-table">
<!-- env-variables -->

## 通用配置

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_APP_KEY** | `""` | App Key (for generation of JWT) | app_key (App Key) |
| **ATK_DEBUG** | `false` | Debug mode | debug (Debug mode) |
| **ATK_HOST** | `"0.0.0.0"` | Listen host | host (Listen host) |
| **ATK_LOCALE** | `"en"` | Language (follow Unicode BCP 47) (可选：`["en", "zh-CN", "zh-TW", "ja", "fr", "ko", "ru"]`) | locale (Language) |
| **ATK_LOGIN_TIMEOUT** | `259200` | Login timeout (in seconds) | login_timeout (Login timeout) |
| **ATK_PORT** | `23366` | Listen port | port (Listen port) |
| **ATK_SITE_DEFAULT** | `"Default Site"` | Default site name (create when app is first launched) | site_default (Default site name) |
| **ATK_SITE_URL** | `""` | Default site url | site_url (Default site url) |
| **ATK_TIMEZONE** | `"Asia/Shanghai"` | Timezone (follow IANA Time Zone Database) | timezone (Timezone) |
| **ATK_TRUSTED_DOMAINS** | `[]` | Trusted domains | trusted_domains (Trusted domains) |


## Multi-Push

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_ADMIN_NOTIFY_BARK_ENABLED** | `false` | 启用 | admin_notify.bark.enabled (Multi-Push > Bark > Enabled) |
| **ATK_ADMIN_NOTIFY_BARK_SERVER** | `"http://day.app/xxxxxxx/"` | Server | admin_notify.bark.server (Multi-Push > Bark > Server) |
| **ATK_ADMIN_NOTIFY_DING_TALK_ENABLED** | `false` | 启用 | admin_notify.ding_talk.enabled (Multi-Push > DingTalk > Enabled) |
| **ATK_ADMIN_NOTIFY_DING_TALK_SECRET** | `""` | Secret | admin_notify.ding_talk.secret (Multi-Push > DingTalk > Secret) |
| **ATK_ADMIN_NOTIFY_DING_TALK_TOKEN** | `""` | Token | admin_notify.ding_talk.token (Multi-Push > DingTalk > Token) |
| **ATK_ADMIN_NOTIFY_EMAIL_ENABLED** | `true` | Enable (can be disabled when using other push methods) | admin_notify.email.enabled (Multi-Push > Notify admin > Enable) |
| **ATK_ADMIN_NOTIFY_EMAIL_MAIL_SUBJECT** | `"[{{site_name}}] Post \"{{page_title}}\" has new a comment"` | Email subject (email subject sent to admin) | admin_notify.email.mail_subject (Multi-Push > Notify admin > Email subject) |
| **ATK_ADMIN_NOTIFY_EMAIL_MAIL_TPL** | `""` | Admin email template file (set to file path to use custom template) | admin_notify.email.mail_tpl (Multi-Push > Notify admin > Admin email template file) |
| **ATK_ADMIN_NOTIFY_LARK_ENABLED** | `false` | 启用 | admin_notify.lark.enabled (Multi-Push > Lark > Enabled) |
| **ATK_ADMIN_NOTIFY_LARK_MSG_TYPE** | `"text"` | Message type (可选：`["text", "card"]`) | admin_notify.lark.msg_type (Multi-Push > Lark > Message type) |
| **ATK_ADMIN_NOTIFY_LARK_WEBHOOK_URL** | `""` | WebhookUrl | admin_notify.lark.webhook_url (Multi-Push > Lark > WebhookUrl) |
| **ATK_ADMIN_NOTIFY_LINE_CHANNEL_ACCESS_TOKEN** | `""` | ChannelAccessToken | admin_notify.line.channel_access_token (Multi-Push > LINE > ChannelAccessToken) |
| **ATK_ADMIN_NOTIFY_LINE_CHANNEL_SECRET** | `""` | ChannelSecret | admin_notify.line.channel_secret (Multi-Push > LINE > ChannelSecret) |
| **ATK_ADMIN_NOTIFY_LINE_ENABLED** | `false` | 启用 | admin_notify.line.enabled (Multi-Push > LINE > Enabled) |
| **ATK_ADMIN_NOTIFY_LINE_RECEIVERS** | `[USER_ID_1 GROUP_ID_1]` | Receivers | admin_notify.line.receivers (Multi-Push > LINE > Receivers) |
| **ATK_ADMIN_NOTIFY_NOISE_MODE** | `false` | Noise mode | admin_notify.noise_mode (Multi-Push > Noise mode) |
| **ATK_ADMIN_NOTIFY_NOTIFY_PENDING** | `false` | Pending comment still send notification (notifications are still sent when comments are intercepted) | admin_notify.notify_pending (Multi-Push > Pending comment still send notification) |
| **ATK_ADMIN_NOTIFY_NOTIFY_TPL** | `"default"` | Notification template (set to file path to use custom template) | admin_notify.notify_tpl (Multi-Push > Notification template) |
| **ATK_ADMIN_NOTIFY_SLACK_ENABLED** | `false` | 启用 | admin_notify.slack.enabled (Multi-Push > Slack > Enabled) |
| **ATK_ADMIN_NOTIFY_SLACK_OAUTH_TOKEN** | `""` | OauthToken | admin_notify.slack.oauth_token (Multi-Push > Slack > OauthToken) |
| **ATK_ADMIN_NOTIFY_SLACK_RECEIVERS** | `[CHANNEL_ID]` | Receivers | admin_notify.slack.receivers (Multi-Push > Slack > Receivers) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_API_TOKEN** | `""` | ApiToken | admin_notify.telegram.api_token (Multi-Push > Telegram > ApiToken) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_ENABLED** | `false` | 启用 | admin_notify.telegram.enabled (Multi-Push > Telegram > Enabled) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_RECEIVERS** | `[7777777]` | Receivers | admin_notify.telegram.receivers (Multi-Push > Telegram > Receivers) |
| **ATK_ADMIN_NOTIFY_WEBHOOK_ENABLED** | `false` | 启用 | admin_notify.webhook.enabled (Multi-Push > WebHook > Enabled) |
| **ATK_ADMIN_NOTIFY_WEBHOOK_URL** | `""` | Url | admin_notify.webhook.url (Multi-Push > WebHook > Url) |


## Social Login

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_AUTH_ANONYMOUS** | `false` | Allow anonymous commenting (Allow skipping verification, only fill in an anonymous nickname and email) | auth.anonymous (Social Login > Allow anonymous commenting) |
| **ATK_AUTH_APPLE_CLIENT_ID** | `""` | ClientId | auth.apple.client_id (Social Login > Apple > ClientId) |
| **ATK_AUTH_APPLE_CLIENT_SECRET** | `""` | ClientSecret | auth.apple.client_secret (Social Login > Apple > ClientSecret) |
| **ATK_AUTH_APPLE_ENABLED** | `false` | 启用 | auth.apple.enabled (Social Login > Apple > Enabled) |
| **ATK_AUTH_AUTH0_CLIENT_ID** | `""` | ClientId | auth.auth0.client_id (Social Login > Auth0 > ClientId) |
| **ATK_AUTH_AUTH0_CLIENT_SECRET** | `""` | ClientSecret | auth.auth0.client_secret (Social Login > Auth0 > ClientSecret) |
| **ATK_AUTH_AUTH0_DOMAIN** | `""` | Domain | auth.auth0.domain (Social Login > Auth0 > Domain) |
| **ATK_AUTH_AUTH0_ENABLED** | `false` | 启用 | auth.auth0.enabled (Social Login > Auth0 > Enabled) |
| **ATK_AUTH_CALLBACK** | `"http://localhost:23366/api/v2/auth/{provider}/callback"` | Callback URL (https://example.com/api/v2/auth/{provider}/callback) | auth.callback (Social Login > Callback URL) |
| **ATK_AUTH_DISCORD_CLIENT_ID** | `""` | ClientId | auth.discord.client_id (Social Login > Discord > ClientId) |
| **ATK_AUTH_DISCORD_CLIENT_SECRET** | `""` | ClientSecret | auth.discord.client_secret (Social Login > Discord > ClientSecret) |
| **ATK_AUTH_DISCORD_ENABLED** | `false` | 启用 | auth.discord.enabled (Social Login > Discord > Enabled) |
| **ATK_AUTH_EMAIL_ENABLED** | `true` | Enable email password login | auth.email.enabled (Social Login > Email > Enable email password login) |
| **ATK_AUTH_EMAIL_VERIFY_SUBJECT** | `"Your Code - {{code}}"` | Verification email subject | auth.email.verify_subject (Social Login > Email > Verification email subject) |
| **ATK_AUTH_EMAIL_VERIFY_TPL** | `"default"` | Verification email template (set to file path to use custom template) | auth.email.verify_tpl (Social Login > Email > Verification email template) |
| **ATK_AUTH_ENABLED** | `false` | Enable Social Login | auth.enabled (Social Login > Enable Social Login) |
| **ATK_AUTH_FACEBOOK_CLIENT_ID** | `""` | ClientId | auth.facebook.client_id (Social Login > Facebook > ClientId) |
| **ATK_AUTH_FACEBOOK_CLIENT_SECRET** | `""` | ClientSecret | auth.facebook.client_secret (Social Login > Facebook > ClientSecret) |
| **ATK_AUTH_FACEBOOK_ENABLED** | `false` | 启用 | auth.facebook.enabled (Social Login > Facebook > Enabled) |
| **ATK_AUTH_GITEA_CLIENT_ID** | `""` | ClientId | auth.gitea.client_id (Social Login > Gitea > ClientId) |
| **ATK_AUTH_GITEA_CLIENT_SECRET** | `""` | ClientSecret | auth.gitea.client_secret (Social Login > Gitea > ClientSecret) |
| **ATK_AUTH_GITEA_ENABLED** | `false` | 启用 | auth.gitea.enabled (Social Login > Gitea > Enabled) |
| **ATK_AUTH_GITHUB_CLIENT_ID** | `""` | ClientId | auth.github.client_id (Social Login > GitHub > ClientId) |
| **ATK_AUTH_GITHUB_CLIENT_SECRET** | `""` | ClientSecret | auth.github.client_secret (Social Login > GitHub > ClientSecret) |
| **ATK_AUTH_GITHUB_ENABLED** | `false` | 启用 | auth.github.enabled (Social Login > GitHub > Enabled) |
| **ATK_AUTH_GITLAB_CLIENT_ID** | `""` | ClientId | auth.gitlab.client_id (Social Login > GitLab > ClientId) |
| **ATK_AUTH_GITLAB_CLIENT_SECRET** | `""` | ClientSecret | auth.gitlab.client_secret (Social Login > GitLab > ClientSecret) |
| **ATK_AUTH_GITLAB_ENABLED** | `false` | 启用 | auth.gitlab.enabled (Social Login > GitLab > Enabled) |
| **ATK_AUTH_GOOGLE_CLIENT_ID** | `""` | ClientId | auth.google.client_id (Social Login > Google > ClientId) |
| **ATK_AUTH_GOOGLE_CLIENT_SECRET** | `""` | ClientSecret | auth.google.client_secret (Social Login > Google > ClientSecret) |
| **ATK_AUTH_GOOGLE_ENABLED** | `false` | 启用 | auth.google.enabled (Social Login > Google > Enabled) |
| **ATK_AUTH_LINE_CLIENT_ID** | `""` | ClientId | auth.line.client_id (Social Login > Line > ClientId) |
| **ATK_AUTH_LINE_CLIENT_SECRET** | `""` | ClientSecret | auth.line.client_secret (Social Login > Line > ClientSecret) |
| **ATK_AUTH_LINE_ENABLED** | `false` | 启用 | auth.line.enabled (Social Login > Line > Enabled) |
| **ATK_AUTH_MASTODON_CLIENT_ID** | `""` | ClientId | auth.mastodon.client_id (Social Login > Mastodon > ClientId) |
| **ATK_AUTH_MASTODON_CLIENT_SECRET** | `""` | ClientSecret | auth.mastodon.client_secret (Social Login > Mastodon > ClientSecret) |
| **ATK_AUTH_MASTODON_ENABLED** | `false` | 启用 | auth.mastodon.enabled (Social Login > Mastodon > Enabled) |
| **ATK_AUTH_MICROSOFT_CLIENT_ID** | `""` | ClientId | auth.microsoft.client_id (Social Login > Microsoft > ClientId) |
| **ATK_AUTH_MICROSOFT_CLIENT_SECRET** | `""` | ClientSecret | auth.microsoft.client_secret (Social Login > Microsoft > ClientSecret) |
| **ATK_AUTH_MICROSOFT_ENABLED** | `false` | 启用 | auth.microsoft.enabled (Social Login > Microsoft > Enabled) |
| **ATK_AUTH_PATREON_CLIENT_ID** | `""` | ClientId | auth.patreon.client_id (Social Login > Patreon > ClientId) |
| **ATK_AUTH_PATREON_CLIENT_SECRET** | `""` | ClientSecret | auth.patreon.client_secret (Social Login > Patreon > ClientSecret) |
| **ATK_AUTH_PATREON_ENABLED** | `false` | 启用 | auth.patreon.enabled (Social Login > Patreon > Enabled) |
| **ATK_AUTH_SLACK_CLIENT_ID** | `""` | ClientId | auth.slack.client_id (Social Login > Slack > ClientId) |
| **ATK_AUTH_SLACK_CLIENT_SECRET** | `""` | ClientSecret | auth.slack.client_secret (Social Login > Slack > ClientSecret) |
| **ATK_AUTH_SLACK_ENABLED** | `false` | 启用 | auth.slack.enabled (Social Login > Slack > Enabled) |
| **ATK_AUTH_STEAM_API_KEY** | `""` | ApiKey | auth.steam.api_key (Social Login > Steam > ApiKey) |
| **ATK_AUTH_STEAM_ENABLED** | `false` | 启用 | auth.steam.enabled (Social Login > Steam > Enabled) |
| **ATK_AUTH_TIKTOK_CLIENT_ID** | `""` | ClientId | auth.tiktok.client_id (Social Login > Tiktok > ClientId) |
| **ATK_AUTH_TIKTOK_CLIENT_SECRET** | `""` | ClientSecret | auth.tiktok.client_secret (Social Login > Tiktok > ClientSecret) |
| **ATK_AUTH_TIKTOK_ENABLED** | `false` | 启用 | auth.tiktok.enabled (Social Login > Tiktok > Enabled) |
| **ATK_AUTH_TWITTER_CLIENT_ID** | `""` | ClientId | auth.twitter.client_id (Social Login > Twitter > ClientId) |
| **ATK_AUTH_TWITTER_CLIENT_SECRET** | `""` | ClientSecret | auth.twitter.client_secret (Social Login > Twitter > ClientSecret) |
| **ATK_AUTH_TWITTER_ENABLED** | `false` | 启用 | auth.twitter.enabled (Social Login > Twitter > Enabled) |
| **ATK_AUTH_WECHAT_CLIENT_ID** | `""` | ClientId | auth.wechat.client_id (Social Login > WeChat > ClientId) |
| **ATK_AUTH_WECHAT_CLIENT_SECRET** | `""` | ClientSecret | auth.wechat.client_secret (Social Login > WeChat > ClientSecret) |
| **ATK_AUTH_WECHAT_ENABLED** | `false` | 启用 | auth.wechat.enabled (Social Login > WeChat > Enabled) |


## Cache

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_CACHE_ENABLED** | `false` | Enable cache | cache.enabled (Cache > Enable cache) |
| **ATK_CACHE_EXPIRES** | `30` | Cache expiration time (in minutes) | cache.expires (Cache > Cache expiration time) |
| **ATK_CACHE_REDIS_DB** | `0` | Redis database number (e.g. 0) | cache.redis.db (Cache > Redis config > Redis database number) |
| **ATK_CACHE_REDIS_NETWORK** | `"tcp"` | Connection type (可选：`["tcp", "unix"]`) | cache.redis.network (Cache > Redis config > Connection type) |
| **ATK_CACHE_REDIS_PASSWORD** | `""` | Redis password | cache.redis.password (Cache > Redis config > Redis password) |
| **ATK_CACHE_REDIS_USERNAME** | `""` | Redis username | cache.redis.username (Cache > Redis config > Redis username) |
| **ATK_CACHE_SERVER** | `""` | Cache server address (e.g. "localhost:6379") | cache.server (Cache > Cache server address) |
| **ATK_CACHE_TYPE** | `"builtin"` | Cache type (可选：`["redis", "memcache", "builtin"]`) | cache.type (Cache > Cache type) |
| **ATK_CACHE_WARM_UP** | `false` | Cache warm up (warm up cache when program starts) | cache.warm_up (Cache > Cache warm up) |


## Captcha

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_CAPTCHA_ACTION_LIMIT** | `3` | Action limit (the number of actions required to activate captcha) | captcha.action_limit (Captcha > Action limit) |
| **ATK_CAPTCHA_ACTION_RESET** | `60` | Reset Timeout (timeout to reset action counter. unit: s, set to -1 to disable) | captcha.action_reset (Captcha > Reset Timeout) |
| **ATK_CAPTCHA_ALWAYS** | `false` | Captcha is required always | captcha.always (Captcha > Captcha is required always) |
| **ATK_CAPTCHA_CAPTCHA_TYPE** | `"image"` | Captcha type (可选：`["image", "turnstile", "recaptcha", "hcaptcha", "geetest"]`) | captcha.captcha_type (Captcha > Captcha type) |
| **ATK_CAPTCHA_ENABLED** | `true` | Enable captcha | captcha.enabled (Captcha > Enable captcha) |
| **ATK_CAPTCHA_GEETEST_CAPTCHA_ID** | `""` | CaptchaId | captcha.geetest.captcha_id (Captcha > Geetest > CaptchaId) |
| **ATK_CAPTCHA_GEETEST_CAPTCHA_KEY** | `""` | CaptchaKey | captcha.geetest.captcha_key (Captcha > Geetest > CaptchaKey) |
| **ATK_CAPTCHA_HCAPTCHA_SECRET_KEY** | `""` | SecretKey | captcha.hcaptcha.secret_key (Captcha > hCaptcha > SecretKey) |
| **ATK_CAPTCHA_HCAPTCHA_SITE_KEY** | `""` | SiteKey | captcha.hcaptcha.site_key (Captcha > hCaptcha > SiteKey) |
| **ATK_CAPTCHA_RECAPTCHA_SECRET_KEY** | `""` | SecretKey | captcha.recaptcha.secret_key (Captcha > reCaptcha > SecretKey) |
| **ATK_CAPTCHA_RECAPTCHA_SITE_KEY** | `""` | SiteKey | captcha.recaptcha.site_key (Captcha > reCaptcha > SiteKey) |
| **ATK_CAPTCHA_TURNSTILE_SECRET_KEY** | `""` | SecretKey | captcha.turnstile.secret_key (Captcha > Turnstile > SecretKey) |
| **ATK_CAPTCHA_TURNSTILE_SITE_KEY** | `""` | SiteKey | captcha.turnstile.site_key (Captcha > Turnstile > SiteKey) |


## Database

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_DB_CHARSET** | `"utf8mb4"` | Database charset | db.charset (Database > Database charset) |
| **ATK_DB_FILE** | `"./data/artalk.db"` | Database file (only for SQLite) | db.file (Database > Database file) |
| **ATK_DB_HOST** | `"localhost"` | Host address | db.host (Database > Host address) |
| **ATK_DB_NAME** | `"artalk"` | Database name | db.name (Database > Database name) |
| **ATK_DB_PASSWORD** | `""` | Database password | db.password (Database > Database password) |
| **ATK_DB_PORT** | `3306` | Host port | db.port (Database > Host port) |
| **ATK_DB_PREPARE_STMT** | `true` | Prepared Statement | db.prepare_stmt (Database > Prepared Statement) |
| **ATK_DB_SSL** | `false` | Enable SSL mode | db.ssl (Database > Enable SSL mode) |
| **ATK_DB_TABLE_PREFIX** | `""` | Table prefix (e.g. "atk_") | db.table_prefix (Database > Table prefix) |
| **ATK_DB_TYPE** | `"sqlite"` | Database type (可选：`["sqlite", "mysql", "pgsql", "mssql"]`) | db.type (Database > Database type) |
| **ATK_DB_USER** | `"root"` | Database user | db.user (Database > Database user) |


## Email

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_EMAIL_ALI_DM_ACCESS_KEY_ID** | `""` | AccessKeyId | email.ali_dm.access_key_id (Email > Aliyun mail push > AccessKeyId) |
| **ATK_EMAIL_ALI_DM_ACCESS_KEY_SECRET** | `""` | AccessKeySecret | email.ali_dm.access_key_secret (Email > Aliyun mail push > AccessKeySecret) |
| **ATK_EMAIL_ALI_DM_ACCOUNT_NAME** | `"noreply@example.com"` | AccountName | email.ali_dm.account_name (Email > Aliyun mail push > AccountName) |
| **ATK_EMAIL_ENABLED** | `false` | Enable email notification | email.enabled (Email > Enable email notification) |
| **ATK_EMAIL_MAIL_SUBJECT** | `"[{{site_name}}] You got a reply from @{{reply_nick}}"` | Email subject | email.mail_subject (Email > Email subject) |
| **ATK_EMAIL_MAIL_TPL** | `"default"` | Email template file (set to file path to use custom template) | email.mail_tpl (Email > Email template file) |
| **ATK_EMAIL_SEND_ADDR** | `"noreply@example.com"` | Email address of sender | email.send_addr (Email > Email address of sender) |
| **ATK_EMAIL_SEND_NAME** | `"{{reply_nick}}"` | Nick name of sender | email.send_name (Email > Nick name of sender) |
| **ATK_EMAIL_SEND_TYPE** | `"smtp"` | Send method (可选：`["smtp", "ali_dm", "sendmail"]`) | email.send_type (Email > Send method) |
| **ATK_EMAIL_SMTP_HOST** | `"smtp.qq.com"` | Email address of sender | email.smtp.host (Email > SMTP send > Email address of sender) |
| **ATK_EMAIL_SMTP_PASSWORD** | `""` | Password | email.smtp.password (Email > SMTP send > Password) |
| **ATK_EMAIL_SMTP_PORT** | `587` | Email port | email.smtp.port (Email > SMTP send > Email port) |
| **ATK_EMAIL_SMTP_USERNAME** | `"example@qq.com"` | Email address of sender | email.smtp.username (Email > SMTP send > Email address of sender) |


## UI Settings

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_FRONTEND_DARKMODE** | `"inherit"` | Dark mode (可选：`["inherit", "auto"]`) | frontend.darkMode (UI Settings > Dark mode) |
| **ATK_FRONTEND_EDITORTRAVEL** | `true` | Movable comment box | frontend.editorTravel (UI Settings > Movable comment box) |
| **ATK_FRONTEND_EMOTICONS** | `"https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json"` | Emoticons | frontend.emoticons (UI Settings > Emoticons) |
| **ATK_FRONTEND_FLATMODE** | `"auto"` | Flatten mode (可选：`["auto", "true", "false"]`) | frontend.flatMode (UI Settings > Flatten mode) |
| **ATK_FRONTEND_GRAVATAR_MIRROR** | `"https://www.gravatar.com/avatar/"` | API URL | frontend.gravatar.mirror (UI Settings > Gravatar > API URL) |
| **ATK_FRONTEND_GRAVATAR_PARAMS** | `"sha256=1&d=mp&s=240"` | API parameters | frontend.gravatar.params (UI Settings > Gravatar > API parameters) |
| **ATK_FRONTEND_HEIGHTLIMIT_CHILDREN** | `400` | Sub-comment area height limit (unit: px) | frontend.heightLimit.children (UI Settings > Content height limit > Sub-comment area height limit) |
| **ATK_FRONTEND_HEIGHTLIMIT_CONTENT** | `300` | Comment content height limit (unit: px) | frontend.heightLimit.content (UI Settings > Content height limit > Comment content height limit) |
| **ATK_FRONTEND_HEIGHTLIMIT_SCROLLABLE** | `false` | Scrollable (scrollable height limit area) | frontend.heightLimit.scrollable (UI Settings > Content height limit > Scrollable) |
| **ATK_FRONTEND_IMGLAZYLOAD** | `false` | Image lazy load (可选：`["false", "native", "data-src"]`) | frontend.imgLazyLoad (UI Settings > Image lazy load) |
| **ATK_FRONTEND_LISTSORT** | `true` | Comment sorting | frontend.listSort (UI Settings > Comment sorting) |
| **ATK_FRONTEND_NESTMAX** | `2` | Maximum nesting level | frontend.nestMax (UI Settings > Maximum nesting level) |
| **ATK_FRONTEND_NESTSORT** | `"DATE_ASC"` | Nesting comment sorting rules (可选：`["DATE_ASC", "DATE_DESC", "VOTE_UP_DESC"]`) | frontend.nestSort (UI Settings > Nesting comment sorting rules) |
| **ATK_FRONTEND_NOCOMMENT** | `""` | Text to display when there is | frontend.noComment (UI Settings > Text to display when there is) |
| **ATK_FRONTEND_PAGINATION_AUTOLOAD** | `true` | Scroll loading | frontend.pagination.autoLoad (UI Settings > Comment pagination > Scroll loading) |
| **ATK_FRONTEND_PAGINATION_PAGESIZE** | `20` | Number of comments per page | frontend.pagination.pageSize (UI Settings > Comment pagination > Number of comments per page) |
| **ATK_FRONTEND_PAGINATION_READMORE** | `true` | Load more mode (disabled to use pagination bar) | frontend.pagination.readMore (UI Settings > Comment pagination > Load more mode) |
| **ATK_FRONTEND_PLACEHOLDER** | `""` | Comment box placeholder | frontend.placeholder (UI Settings > Comment box placeholder) |
| **ATK_FRONTEND_PLUGINURLS** | `[]` | Plugins | frontend.pluginURLs (UI Settings > Plugins) |
| **ATK_FRONTEND_PREVIEW** | `true` | Editor real-time preview | frontend.preview (UI Settings > Editor real-time preview) |
| **ATK_FRONTEND_REQTIMEOUT** | `15000` | Request timeout (unit: ms) | frontend.reqTimeout (UI Settings > Request timeout) |
| **ATK_FRONTEND_SENDBTN** | `""` | Text of the send button | frontend.sendBtn (UI Settings > Text of the send button) |
| **ATK_FRONTEND_UABADGE** | `false` | User UA badge | frontend.uaBadge (UI Settings > User UA badge) |
| **ATK_FRONTEND_VERSIONCHECK** | `true` | Version check | frontend.versionCheck (UI Settings > Version check) |
| **ATK_FRONTEND_VOTE** | `true` | Vote button | frontend.vote (UI Settings > Vote button) |
| **ATK_FRONTEND_VOTEDOWN** | `false` | Vote down button | frontend.voteDown (UI Settings > Vote down button) |


## Web server

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_HTTP_BODY_LIMIT** | `100` | Body size limit (unit: MB) | http.body_limit (Web server > Body size limit) |
| **ATK_HTTP_PROXY_HEADER** | `""` | Proxy Header (fill `X-Forwarded-For` to get user real IP if behind a trusted reverse proxy or CDN) | http.proxy_header (Web server > Proxy Header) |


## Upload

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_IMG_UPLOAD_ENABLED** | `true` | Enable image upload | img_upload.enabled (Upload > Enable image upload) |
| **ATK_IMG_UPLOAD_MAX_SIZE** | `5` | Image size limit (unit: MB) | img_upload.max_size (Upload > Image size limit) |
| **ATK_IMG_UPLOAD_PATH** | `"./data/artalk-img/"` | Image storage | img_upload.path (Upload > Image storage) |
| **ATK_IMG_UPLOAD_PUBLIC_PATH** | `<nil>` | Image link base path (default: "/static/images/") | img_upload.public_path (Upload > Image link base path) |
| **ATK_IMG_UPLOAD_UPGIT_DEL_LOCAL** | `true` | Delete local image after upload success | img_upload.upgit.del_local (Upload > Upgit config > Delete local image after upload success) |
| **ATK_IMG_UPLOAD_UPGIT_ENABLED** | `false` | Enable Upgit | img_upload.upgit.enabled (Upload > Upgit config > Enable Upgit) |
| **ATK_IMG_UPLOAD_UPGIT_EXEC** | `"upgit -c UPGIT_CONF_FILE_PATH -t /artalk-img"` | Command line arguments | img_upload.upgit.exec (Upload > Upgit config > Command line arguments) |


## Logging

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_LOG_ENABLED** | `true` | Enable logging | log.enabled (Logging > Enable logging) |
| **ATK_LOG_FILENAME** | `"./data/artalk.log"` | Log file path | log.filename (Logging > Log file path) |


## Moderator

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_MODERATOR_AKISMET_KEY** | `""` | Akismet Key (Akismet anti-spam service, https://akismet.com) | moderator.akismet_key (Moderator > Akismet Key) |
| **ATK_MODERATOR_ALIYUN_ACCESS_KEY_ID** | `""` | AccessKeyId | moderator.aliyun.access_key_id (Moderator > Aliyun Content Security > AccessKeyId) |
| **ATK_MODERATOR_ALIYUN_ACCESS_KEY_SECRET** | `""` | AccessKeySecret | moderator.aliyun.access_key_secret (Moderator > Aliyun Content Security > AccessKeySecret) |
| **ATK_MODERATOR_ALIYUN_ENABLED** | `false` | 启用 | moderator.aliyun.enabled (Moderator > Aliyun Content Security > Enabled) |
| **ATK_MODERATOR_ALIYUN_REGION** | `"cn-shanghai"` | Region | moderator.aliyun.region (Moderator > Aliyun Content Security > Region) |
| **ATK_MODERATOR_API_FAIL_BLOCK** | `false` | Block when API request fails (set to false to let comments pass when API request fails) | moderator.api_fail_block (Moderator > Block when API request fails) |
| **ATK_MODERATOR_KEYWORDS_ENABLED** | `false` | Enable keyword filter | moderator.keywords.enabled (Moderator > Keyword filter > Enable keyword filter) |
| **ATK_MODERATOR_KEYWORDS_FILE_SEP** | `"\n"` | FileSep | moderator.keywords.file_sep (Moderator > Keyword filter > FileSep) |
| **ATK_MODERATOR_KEYWORDS_FILES** | `[./data/keywords_1.txt]` | Dictionary file (support multiple dictionary files) | moderator.keywords.files (Moderator > Keyword filter > Dictionary file) |
| **ATK_MODERATOR_KEYWORDS_PENDING** | `false` | Set to pending when match | moderator.keywords.pending (Moderator > Keyword filter > Set to pending when match) |
| **ATK_MODERATOR_KEYWORDS_REPLACE_TO** | `"x"` | ReplaceTo | moderator.keywords.replace_to (Moderator > Keyword filter > ReplaceTo) |
| **ATK_MODERATOR_PENDING_DEFAULT** | `false` | Default pending (new comments need to be approved by admin) | moderator.pending_default (Moderator > Default pending) |
| **ATK_MODERATOR_TENCENT_ENABLED** | `false` | 启用 | moderator.tencent.enabled (Moderator > Tencent Cloud Content Security > Enabled) |
| **ATK_MODERATOR_TENCENT_REGION** | `"ap-guangzhou"` | Region | moderator.tencent.region (Moderator > Tencent Cloud Content Security > Region) |
| **ATK_MODERATOR_TENCENT_SECRET_ID** | `""` | SecretId | moderator.tencent.secret_id (Moderator > Tencent Cloud Content Security > SecretId) |
| **ATK_MODERATOR_TENCENT_SECRET_KEY** | `""` | SecretKey | moderator.tencent.secret_key (Moderator > Tencent Cloud Content Security > SecretKey) |


## SSL

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_SSL_CERT_PATH** | `""` | Certificate file path (e.g. "/etc/letsencrypt/live/example.com/fullchain.pem") | ssl.cert_path (SSL > Certificate file path) |
| **ATK_SSL_ENABLED** | `false` | Enable SSL | ssl.enabled (SSL > Enable SSL) |
| **ATK_SSL_KEY_PATH** | `""` | Key file path (e.g. "/etc/letsencrypt/live/example.com/privkey.pem") | ssl.key_path (SSL > Key file path) |

<!-- /env-variables -->
</div>

## Developers

Enable developer mode by setting `ATK_DEBUG=1`.

This document is auto-generated by the Artalk developer tool. Run `make update-conf-docs` to update this page's content based on the template configuration file ([conf/artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/tree/master/conf/artalk.example.zh-CN.yml)).

Related code can be found in the [@ArtalkJS/Artalk:/internal/config/](https://github.com/ArtalkJS/Artalk/tree/master/internal/config/) directory.
