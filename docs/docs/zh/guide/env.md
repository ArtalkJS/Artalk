# 环境变量

Artalk 支持通过环境变量修改配置，这使得 Docker 等容器化云部署变得简单。

环境变量名以 `ATK_` 开头，全大写，对应[配置文件](./backend/config.md)中的每个节点。

环境变量的优先级高于配置文件，通过以下命令验证配置是否生效：
  
```bash
artalk config
```

使用 Docker Compose 部署 Artalk，可在 `compose.yml` 文件添加环境变量，例如：

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
      - TZ=Asia/Shanghai
      - ATK_LOCALE=zh-CN
      - ATK_SITE_DEFAULT=Artalk 的博客
      - ATK_SITE_URL=https://example.com
      - ATK_TRUSTED_DOMAINS=https://dev.example.com https://localhost:8080
      - ATK_ADMIN_USERS_0_NAME=admin
      - ATK_ADMIN_USERS_0_EMAIL=admin@example.org
      - ATK_ADMIN_USERS_0_PASSWORD=(bcrypt)$2y$10$ti4vZYIrxVN8rLcYXVgXCO.GJND0dyI49r7IoF3xqIx8bBRmIBZRm
      - ATK_ADMIN_USERS_0_BADGE_NAME=管理员
      - ATK_ADMIN_USERS_0_BADGE_COLOR=#0083FF
```

当变量为数组，通过空格分隔的字符串或数字索引来设置数组值，例如：

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
| **ATK_ADMIN_USERS** | `[]` | 管理员账户 | admin_users (管理员账户) |
| **ATK_APP_KEY** | `""` | 加密密钥 | app_key (加密密钥) |
| **ATK_DEBUG** | `false` | 调试模式 | debug (调试模式) |
| **ATK_HOST** | `"0.0.0.0"` | 服务器地址 | host (服务器地址) |
| **ATK_LOCALE** | `"zh-CN"` | 语言 (可选：`["en", "zh-CN", "zh-TW", "ja", "fr", "ko", "ru"]`) | locale (语言) |
| **ATK_LOGIN_TIMEOUT** | `259200` | 登录有效时长 (单位：秒) | login_timeout (登录有效时长) |
| **ATK_PORT** | `23366` | 服务器端口 | port (服务器端口) |
| **ATK_SITE_DEFAULT** | `"默认站点"` | 默认站点名 | site_default (默认站点名) |
| **ATK_SITE_URL** | `""` | 默认站点地址 | site_url (默认站点地址) |
| **ATK_TIMEZONE** | `"Asia/Shanghai"` | 时间区域 | timezone (时间区域) |
| **ATK_TRUSTED_DOMAINS** | `[]` | 可信域名 | trusted_domains (可信域名) |


## 多元推送

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_ADMIN_NOTIFY_BARK_ENABLED** | `false` | 启用 | admin_notify.bark.enabled (多元推送 > Bark > Enabled) |
| **ATK_ADMIN_NOTIFY_BARK_SERVER** | `"http://day.app/xxxxxxx/"` | Server | admin_notify.bark.server (多元推送 > Bark > Server) |
| **ATK_ADMIN_NOTIFY_DING_TALK_ENABLED** | `false` | 启用 | admin_notify.ding_talk.enabled (多元推送 > 钉钉 > Enabled) |
| **ATK_ADMIN_NOTIFY_DING_TALK_SECRET** | `""` | Secret | admin_notify.ding_talk.secret (多元推送 > 钉钉 > Secret) |
| **ATK_ADMIN_NOTIFY_DING_TALK_TOKEN** | `""` | Token | admin_notify.ding_talk.token (多元推送 > 钉钉 > Token) |
| **ATK_ADMIN_NOTIFY_EMAIL_ENABLED** | `true` | 开启 (当使用其他推送方式时，可以关闭管理员邮件通知) | admin_notify.email.enabled (多元推送 > 邮件通知管理员 > 开启) |
| **ATK_ADMIN_NOTIFY_EMAIL_MAIL_SUBJECT** | `"[{{site_name}}] 您的文章「{{page_title}}」有新回复"` | 邮件标题 (发送给管理员的邮件标题) | admin_notify.email.mail_subject (多元推送 > 邮件通知管理员 > 邮件标题) |
| **ATK_ADMIN_NOTIFY_EMAIL_MAIL_TPL** | `""` | 管理员邮件模板文件 (填入文件路径使用自定义模板) | admin_notify.email.mail_tpl (多元推送 > 邮件通知管理员 > 管理员邮件模板文件) |
| **ATK_ADMIN_NOTIFY_LARK_ENABLED** | `false` | 启用 | admin_notify.lark.enabled (多元推送 > 飞书 > Enabled) |
| **ATK_ADMIN_NOTIFY_LARK_MSG_TYPE** | `"text"` | 消息类型 (可选：`["text", "card"]`) | admin_notify.lark.msg_type (多元推送 > 飞书 > 消息类型) |
| **ATK_ADMIN_NOTIFY_LARK_WEBHOOK_URL** | `""` | WebhookUrl | admin_notify.lark.webhook_url (多元推送 > 飞书 > WebhookUrl) |
| **ATK_ADMIN_NOTIFY_LINE_CHANNEL_ACCESS_TOKEN** | `""` | ChannelAccessToken | admin_notify.line.channel_access_token (多元推送 > LINE > ChannelAccessToken) |
| **ATK_ADMIN_NOTIFY_LINE_CHANNEL_SECRET** | `""` | ChannelSecret | admin_notify.line.channel_secret (多元推送 > LINE > ChannelSecret) |
| **ATK_ADMIN_NOTIFY_LINE_ENABLED** | `false` | 启用 | admin_notify.line.enabled (多元推送 > LINE > Enabled) |
| **ATK_ADMIN_NOTIFY_LINE_RECEIVERS** | `[USER_ID_1 GROUP_ID_1]` | Receivers | admin_notify.line.receivers (多元推送 > LINE > Receivers) |
| **ATK_ADMIN_NOTIFY_NOISE_MODE** | `false` | 嘈杂模式 | admin_notify.noise_mode (多元推送 > 嘈杂模式) |
| **ATK_ADMIN_NOTIFY_NOTIFY_PENDING** | `false` | 待审评论仍然发送通知 (当评论被拦截时仍然发送通知) | admin_notify.notify_pending (多元推送 > 待审评论仍然发送通知) |
| **ATK_ADMIN_NOTIFY_NOTIFY_TPL** | `"default"` | 通知模版 (填入文件路径使用自定义模板) | admin_notify.notify_tpl (多元推送 > 通知模版) |
| **ATK_ADMIN_NOTIFY_SLACK_ENABLED** | `false` | 启用 | admin_notify.slack.enabled (多元推送 > Slack > Enabled) |
| **ATK_ADMIN_NOTIFY_SLACK_OAUTH_TOKEN** | `""` | OauthToken | admin_notify.slack.oauth_token (多元推送 > Slack > OauthToken) |
| **ATK_ADMIN_NOTIFY_SLACK_RECEIVERS** | `[CHANNEL_ID]` | Receivers | admin_notify.slack.receivers (多元推送 > Slack > Receivers) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_API_TOKEN** | `""` | ApiToken | admin_notify.telegram.api_token (多元推送 > Telegram > ApiToken) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_ENABLED** | `false` | 启用 | admin_notify.telegram.enabled (多元推送 > Telegram > Enabled) |
| **ATK_ADMIN_NOTIFY_TELEGRAM_RECEIVERS** | `[7777777]` | Receivers | admin_notify.telegram.receivers (多元推送 > Telegram > Receivers) |
| **ATK_ADMIN_NOTIFY_WEBHOOK_ENABLED** | `false` | 启用 | admin_notify.webhook.enabled (多元推送 > WebHook > Enabled) |
| **ATK_ADMIN_NOTIFY_WEBHOOK_URL** | `""` | Url | admin_notify.webhook.url (多元推送 > WebHook > Url) |


## 社交登录

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_AUTH_ANONYMOUS** | `false` | 允许匿名评论 (允许跳过验证，仅填写匿名的昵称和邮箱) | auth.anonymous (社交登录 > 允许匿名评论) |
| **ATK_AUTH_APPLE_CLIENT_ID** | `""` | ClientId | auth.apple.client_id (社交登录 > Apple > ClientId) |
| **ATK_AUTH_APPLE_CLIENT_SECRET** | `""` | ClientSecret | auth.apple.client_secret (社交登录 > Apple > ClientSecret) |
| **ATK_AUTH_APPLE_ENABLED** | `false` | 启用 | auth.apple.enabled (社交登录 > Apple > Enabled) |
| **ATK_AUTH_AUTH0_CLIENT_ID** | `""` | ClientId | auth.auth0.client_id (社交登录 > Auth0 > ClientId) |
| **ATK_AUTH_AUTH0_CLIENT_SECRET** | `""` | ClientSecret | auth.auth0.client_secret (社交登录 > Auth0 > ClientSecret) |
| **ATK_AUTH_AUTH0_DOMAIN** | `""` | Domain | auth.auth0.domain (社交登录 > Auth0 > Domain) |
| **ATK_AUTH_AUTH0_ENABLED** | `false` | 启用 | auth.auth0.enabled (社交登录 > Auth0 > Enabled) |
| **ATK_AUTH_CALLBACK** | `"http://localhost:23366/api/v2/auth/{provider}/callback"` | 回调地址 (https://example.com/api/v2/auth/{provider}/callback) | auth.callback (社交登录 > 回调地址) |
| **ATK_AUTH_DISCORD_CLIENT_ID** | `""` | ClientId | auth.discord.client_id (社交登录 > Discord > ClientId) |
| **ATK_AUTH_DISCORD_CLIENT_SECRET** | `""` | ClientSecret | auth.discord.client_secret (社交登录 > Discord > ClientSecret) |
| **ATK_AUTH_DISCORD_ENABLED** | `false` | 启用 | auth.discord.enabled (社交登录 > Discord > Enabled) |
| **ATK_AUTH_EMAIL_ENABLED** | `true` | 启用邮箱密码登录 | auth.email.enabled (社交登录 > Email > 启用邮箱密码登录) |
| **ATK_AUTH_EMAIL_VERIFY_SUBJECT** | `"您的验证码是 - {{code}}"` | 邮箱验证邮件标题 | auth.email.verify_subject (社交登录 > Email > 邮箱验证邮件标题) |
| **ATK_AUTH_EMAIL_VERIFY_TPL** | `"default"` | 邮箱验证邮件模板 (填入文件路径使用自定义模板) | auth.email.verify_tpl (社交登录 > Email > 邮箱验证邮件模板) |
| **ATK_AUTH_ENABLED** | `false` | 启用社交登录 | auth.enabled (社交登录 > 启用社交登录) |
| **ATK_AUTH_FACEBOOK_CLIENT_ID** | `""` | ClientId | auth.facebook.client_id (社交登录 > Facebook > ClientId) |
| **ATK_AUTH_FACEBOOK_CLIENT_SECRET** | `""` | ClientSecret | auth.facebook.client_secret (社交登录 > Facebook > ClientSecret) |
| **ATK_AUTH_FACEBOOK_ENABLED** | `false` | 启用 | auth.facebook.enabled (社交登录 > Facebook > Enabled) |
| **ATK_AUTH_GITEA_CLIENT_ID** | `""` | ClientId | auth.gitea.client_id (社交登录 > Gitea > ClientId) |
| **ATK_AUTH_GITEA_CLIENT_SECRET** | `""` | ClientSecret | auth.gitea.client_secret (社交登录 > Gitea > ClientSecret) |
| **ATK_AUTH_GITEA_ENABLED** | `false` | 启用 | auth.gitea.enabled (社交登录 > Gitea > Enabled) |
| **ATK_AUTH_GITHUB_CLIENT_ID** | `""` | ClientId | auth.github.client_id (社交登录 > GitHub > ClientId) |
| **ATK_AUTH_GITHUB_CLIENT_SECRET** | `""` | ClientSecret | auth.github.client_secret (社交登录 > GitHub > ClientSecret) |
| **ATK_AUTH_GITHUB_ENABLED** | `false` | 启用 | auth.github.enabled (社交登录 > GitHub > Enabled) |
| **ATK_AUTH_GITLAB_CLIENT_ID** | `""` | ClientId | auth.gitlab.client_id (社交登录 > GitLab > ClientId) |
| **ATK_AUTH_GITLAB_CLIENT_SECRET** | `""` | ClientSecret | auth.gitlab.client_secret (社交登录 > GitLab > ClientSecret) |
| **ATK_AUTH_GITLAB_ENABLED** | `false` | 启用 | auth.gitlab.enabled (社交登录 > GitLab > Enabled) |
| **ATK_AUTH_GOOGLE_CLIENT_ID** | `""` | ClientId | auth.google.client_id (社交登录 > Google > ClientId) |
| **ATK_AUTH_GOOGLE_CLIENT_SECRET** | `""` | ClientSecret | auth.google.client_secret (社交登录 > Google > ClientSecret) |
| **ATK_AUTH_GOOGLE_ENABLED** | `false` | 启用 | auth.google.enabled (社交登录 > Google > Enabled) |
| **ATK_AUTH_LINE_CLIENT_ID** | `""` | ClientId | auth.line.client_id (社交登录 > Line > ClientId) |
| **ATK_AUTH_LINE_CLIENT_SECRET** | `""` | ClientSecret | auth.line.client_secret (社交登录 > Line > ClientSecret) |
| **ATK_AUTH_LINE_ENABLED** | `false` | 启用 | auth.line.enabled (社交登录 > Line > Enabled) |
| **ATK_AUTH_MASTODON_CLIENT_ID** | `""` | ClientId | auth.mastodon.client_id (社交登录 > Mastodon > ClientId) |
| **ATK_AUTH_MASTODON_CLIENT_SECRET** | `""` | ClientSecret | auth.mastodon.client_secret (社交登录 > Mastodon > ClientSecret) |
| **ATK_AUTH_MASTODON_ENABLED** | `false` | 启用 | auth.mastodon.enabled (社交登录 > Mastodon > Enabled) |
| **ATK_AUTH_MICROSOFT_CLIENT_ID** | `""` | ClientId | auth.microsoft.client_id (社交登录 > Microsoft > ClientId) |
| **ATK_AUTH_MICROSOFT_CLIENT_SECRET** | `""` | ClientSecret | auth.microsoft.client_secret (社交登录 > Microsoft > ClientSecret) |
| **ATK_AUTH_MICROSOFT_ENABLED** | `false` | 启用 | auth.microsoft.enabled (社交登录 > Microsoft > Enabled) |
| **ATK_AUTH_PATREON_CLIENT_ID** | `""` | ClientId | auth.patreon.client_id (社交登录 > Patreon > ClientId) |
| **ATK_AUTH_PATREON_CLIENT_SECRET** | `""` | ClientSecret | auth.patreon.client_secret (社交登录 > Patreon > ClientSecret) |
| **ATK_AUTH_PATREON_ENABLED** | `false` | 启用 | auth.patreon.enabled (社交登录 > Patreon > Enabled) |
| **ATK_AUTH_SLACK_CLIENT_ID** | `""` | ClientId | auth.slack.client_id (社交登录 > Slack > ClientId) |
| **ATK_AUTH_SLACK_CLIENT_SECRET** | `""` | ClientSecret | auth.slack.client_secret (社交登录 > Slack > ClientSecret) |
| **ATK_AUTH_SLACK_ENABLED** | `false` | 启用 | auth.slack.enabled (社交登录 > Slack > Enabled) |
| **ATK_AUTH_STEAM_API_KEY** | `""` | ApiKey | auth.steam.api_key (社交登录 > Steam > ApiKey) |
| **ATK_AUTH_STEAM_ENABLED** | `false` | 启用 | auth.steam.enabled (社交登录 > Steam > Enabled) |
| **ATK_AUTH_TIKTOK_CLIENT_ID** | `""` | ClientId | auth.tiktok.client_id (社交登录 > Tiktok > ClientId) |
| **ATK_AUTH_TIKTOK_CLIENT_SECRET** | `""` | ClientSecret | auth.tiktok.client_secret (社交登录 > Tiktok > ClientSecret) |
| **ATK_AUTH_TIKTOK_ENABLED** | `false` | 启用 | auth.tiktok.enabled (社交登录 > Tiktok > Enabled) |
| **ATK_AUTH_TWITTER_CLIENT_ID** | `""` | ClientId | auth.twitter.client_id (社交登录 > Twitter > ClientId) |
| **ATK_AUTH_TWITTER_CLIENT_SECRET** | `""` | ClientSecret | auth.twitter.client_secret (社交登录 > Twitter > ClientSecret) |
| **ATK_AUTH_TWITTER_ENABLED** | `false` | 启用 | auth.twitter.enabled (社交登录 > Twitter > Enabled) |
| **ATK_AUTH_WECHAT_CLIENT_ID** | `""` | ClientId | auth.wechat.client_id (社交登录 > 微信 > ClientId) |
| **ATK_AUTH_WECHAT_CLIENT_SECRET** | `""` | ClientSecret | auth.wechat.client_secret (社交登录 > 微信 > ClientSecret) |
| **ATK_AUTH_WECHAT_ENABLED** | `false` | 启用 | auth.wechat.enabled (社交登录 > 微信 > Enabled) |


## 缓存

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_CACHE_ENABLED** | `false` | 启用缓存 | cache.enabled (缓存 > 启用缓存) |
| **ATK_CACHE_EXPIRES** | `30` | 缓存过期时间 (单位：分钟) | cache.expires (缓存 > 缓存过期时间) |
| **ATK_CACHE_REDIS_DB** | `0` | 数据库编号 (例如使用零号数据库填写 0) | cache.redis.db (缓存 > Redis 配置 > 数据库编号) |
| **ATK_CACHE_REDIS_NETWORK** | `"tcp"` | 连接方式 (可选：`["tcp", "unix"]`) | cache.redis.network (缓存 > Redis 配置 > 连接方式) |
| **ATK_CACHE_REDIS_PASSWORD** | `""` | 密码 | cache.redis.password (缓存 > Redis 配置 > 密码) |
| **ATK_CACHE_REDIS_USERNAME** | `""` | 用户名 | cache.redis.username (缓存 > Redis 配置 > 用户名) |
| **ATK_CACHE_SERVER** | `""` | 缓存服务器地址 (例如："localhost:6379") | cache.server (缓存 > 缓存服务器地址) |
| **ATK_CACHE_TYPE** | `"builtin"` | 缓存类型 (可选：`["redis", "memcache", "builtin"]`) | cache.type (缓存 > 缓存类型) |
| **ATK_CACHE_WARM_UP** | `false` | 缓存启动预热 (程序启动时预热缓存) | cache.warm_up (缓存 > 缓存启动预热) |


## 验证码

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_CAPTCHA_ACTION_LIMIT** | `3` | 激活验证码所需操作次数 | captcha.action_limit (验证码 > 激活验证码所需操作次数) |
| **ATK_CAPTCHA_ACTION_RESET** | `60` | 重置操作计数器超时 (单位：s, 设为 -1 不重置) | captcha.action_reset (验证码 > 重置操作计数器超时) |
| **ATK_CAPTCHA_ALWAYS** | `false` | 总是需要验证码 | captcha.always (验证码 > 总是需要验证码) |
| **ATK_CAPTCHA_CAPTCHA_TYPE** | `"image"` | 验证类型 (可选：`["image", "turnstile", "recaptcha", "hcaptcha", "geetest"]`) | captcha.captcha_type (验证码 > 验证类型) |
| **ATK_CAPTCHA_ENABLED** | `true` | 启用验证码 | captcha.enabled (验证码 > 启用验证码) |
| **ATK_CAPTCHA_GEETEST_CAPTCHA_ID** | `""` | CaptchaId | captcha.geetest.captcha_id (验证码 > Geetest 极验 > CaptchaId) |
| **ATK_CAPTCHA_GEETEST_CAPTCHA_KEY** | `""` | CaptchaKey | captcha.geetest.captcha_key (验证码 > Geetest 极验 > CaptchaKey) |
| **ATK_CAPTCHA_HCAPTCHA_SECRET_KEY** | `""` | SecretKey | captcha.hcaptcha.secret_key (验证码 > hCaptcha > SecretKey) |
| **ATK_CAPTCHA_HCAPTCHA_SITE_KEY** | `""` | SiteKey | captcha.hcaptcha.site_key (验证码 > hCaptcha > SiteKey) |
| **ATK_CAPTCHA_RECAPTCHA_SECRET_KEY** | `""` | SecretKey | captcha.recaptcha.secret_key (验证码 > reCaptcha > SecretKey) |
| **ATK_CAPTCHA_RECAPTCHA_SITE_KEY** | `""` | SiteKey | captcha.recaptcha.site_key (验证码 > reCaptcha > SiteKey) |
| **ATK_CAPTCHA_TURNSTILE_SECRET_KEY** | `""` | SecretKey | captcha.turnstile.secret_key (验证码 > Turnstile > SecretKey) |
| **ATK_CAPTCHA_TURNSTILE_SITE_KEY** | `""` | SiteKey | captcha.turnstile.site_key (验证码 > Turnstile > SiteKey) |


## 数据库

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_DB_CHARSET** | `"utf8mb4"` | 编码格式 | db.charset (数据库 > 编码格式) |
| **ATK_DB_FILE** | `"./data/artalk.db"` | 数据库文件 (仅 SQLite 数据库需填写) | db.file (数据库 > 数据库文件) |
| **ATK_DB_HOST** | `"localhost"` | 数据库地址 | db.host (数据库 > 数据库地址) |
| **ATK_DB_NAME** | `"artalk"` | 数据库名称 | db.name (数据库 > 数据库名称) |
| **ATK_DB_PASSWORD** | `""` | 数据库密码 | db.password (数据库 > 数据库密码) |
| **ATK_DB_PORT** | `3306` | 数据库端口 | db.port (数据库 > 数据库端口) |
| **ATK_DB_PREPARE_STMT** | `true` | 预编译语句 | db.prepare_stmt (数据库 > 预编译语句) |
| **ATK_DB_SSL** | `false` | 启用 SSL | db.ssl (数据库 > 启用 SSL) |
| **ATK_DB_TABLE_PREFIX** | `""` | 表前缀 (例如："atk_") | db.table_prefix (数据库 > 表前缀) |
| **ATK_DB_TYPE** | `"sqlite"` | 数据库类型 (可选：`["sqlite", "mysql", "pgsql", "mssql"]`) | db.type (数据库 > 数据库类型) |
| **ATK_DB_USER** | `"root"` | 数据库账户 | db.user (数据库 > 数据库账户) |


## 邮件通知

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_EMAIL_ALI_DM_ACCESS_KEY_ID** | `""` | AccessKeyId | email.ali_dm.access_key_id (邮件通知 > 阿里云邮件推送 > AccessKeyId) |
| **ATK_EMAIL_ALI_DM_ACCESS_KEY_SECRET** | `""` | AccessKeySecret | email.ali_dm.access_key_secret (邮件通知 > 阿里云邮件推送 > AccessKeySecret) |
| **ATK_EMAIL_ALI_DM_ACCOUNT_NAME** | `"noreply@example.com"` | AccountName | email.ali_dm.account_name (邮件通知 > 阿里云邮件推送 > AccountName) |
| **ATK_EMAIL_ENABLED** | `false` | 启用邮件通知 | email.enabled (邮件通知 > 启用邮件通知) |
| **ATK_EMAIL_MAIL_SUBJECT** | `"[{{site_name}}] 您收到了来自 @{{reply_nick}} 的回复"` | 邮件标题 | email.mail_subject (邮件通知 > 邮件标题) |
| **ATK_EMAIL_MAIL_TPL** | `"default"` | 邮件模板文件 (填入文件路径使用自定义模板) | email.mail_tpl (邮件通知 > 邮件模板文件) |
| **ATK_EMAIL_SEND_ADDR** | `"noreply@example.com"` | 发信人地址 | email.send_addr (邮件通知 > 发信人地址) |
| **ATK_EMAIL_SEND_NAME** | `"{{reply_nick}}"` | 发信人昵称 | email.send_name (邮件通知 > 发信人昵称) |
| **ATK_EMAIL_SEND_TYPE** | `"smtp"` | 发送方式 (可选：`["smtp", "ali_dm", "sendmail"]`) | email.send_type (邮件通知 > 发送方式) |
| **ATK_EMAIL_SMTP_HOST** | `"smtp.qq.com"` | 发件地址 | email.smtp.host (邮件通知 > SMTP 发送 > 发件地址) |
| **ATK_EMAIL_SMTP_PASSWORD** | `""` | 密码 | email.smtp.password (邮件通知 > SMTP 发送 > 密码) |
| **ATK_EMAIL_SMTP_PORT** | `587` | 发件端口 | email.smtp.port (邮件通知 > SMTP 发送 > 发件端口) |
| **ATK_EMAIL_SMTP_USERNAME** | `"example@qq.com"` | 用户名 | email.smtp.username (邮件通知 > SMTP 发送 > 用户名) |


## 界面配置

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_FRONTEND_DARKMODE** | `"inherit"` | 夜间模式 (可选：`["inherit", "auto"]`) | frontend.darkMode (界面配置 > 夜间模式) |
| **ATK_FRONTEND_EDITORTRAVEL** | `true` | 评论框穿梭 | frontend.editorTravel (界面配置 > 评论框穿梭) |
| **ATK_FRONTEND_EMOTICONS** | `"https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json"` | 表情包 | frontend.emoticons (界面配置 > 表情包) |
| **ATK_FRONTEND_FLATMODE** | `"auto"` | 平铺模式 (可选：`["auto", "true", "false"]`) | frontend.flatMode (界面配置 > 平铺模式) |
| **ATK_FRONTEND_GRAVATAR_MIRROR** | `"https://weavatar.com/avatar/"` | API 地址 | frontend.gravatar.mirror (界面配置 > 头像 Gravatar > API 地址) |
| **ATK_FRONTEND_GRAVATAR_PARAMS** | `"sha256=1&d=mp&s=240"` | API 参数 | frontend.gravatar.params (界面配置 > 头像 Gravatar > API 参数) |
| **ATK_FRONTEND_HEIGHTLIMIT_CHILDREN** | `400` | 子评论区域限高 (单位：px) | frontend.heightLimit.children (界面配置 > 内容限高 > 子评论区域限高) |
| **ATK_FRONTEND_HEIGHTLIMIT_CONTENT** | `300` | 评论内容限高 (单位：px) | frontend.heightLimit.content (界面配置 > 内容限高 > 评论内容限高) |
| **ATK_FRONTEND_HEIGHTLIMIT_SCROLLABLE** | `false` | 滚动限高 (允许限高区域滚动) | frontend.heightLimit.scrollable (界面配置 > 内容限高 > 滚动限高) |
| **ATK_FRONTEND_IMGLAZYLOAD** | `false` | 图片懒加载 (可选：`["false", "native", "data-src"]`) | frontend.imgLazyLoad (界面配置 > 图片懒加载) |
| **ATK_FRONTEND_LISTSORT** | `true` | 评论排序功能 | frontend.listSort (界面配置 > 评论排序功能) |
| **ATK_FRONTEND_NESTMAX** | `2` | 最大嵌套层数 | frontend.nestMax (界面配置 > 最大嵌套层数) |
| **ATK_FRONTEND_NESTSORT** | `"DATE_ASC"` | 嵌套评论排序规则 (可选：`["DATE_ASC", "DATE_DESC", "VOTE_UP_DESC"]`) | frontend.nestSort (界面配置 > 嵌套评论排序规则) |
| **ATK_FRONTEND_NOCOMMENT** | `""` | 无评论显示文字 | frontend.noComment (界面配置 > 无评论显示文字) |
| **ATK_FRONTEND_PAGINATION_AUTOLOAD** | `true` | 滚动加载 | frontend.pagination.autoLoad (界面配置 > 评论分页 > 滚动加载) |
| **ATK_FRONTEND_PAGINATION_PAGESIZE** | `20` | 每页评论数 | frontend.pagination.pageSize (界面配置 > 评论分页 > 每页评论数) |
| **ATK_FRONTEND_PAGINATION_READMORE** | `true` | 加载更多模式 (关闭则使用分页条) | frontend.pagination.readMore (界面配置 > 评论分页 > 加载更多模式) |
| **ATK_FRONTEND_PLACEHOLDER** | `""` | 评论框占位文字 | frontend.placeholder (界面配置 > 评论框占位文字) |
| **ATK_FRONTEND_PLUGINURLS** | `[]` | 插件 | frontend.pluginURLs (界面配置 > 插件) |
| **ATK_FRONTEND_PREVIEW** | `true` | 编辑器实时预览功能 | frontend.preview (界面配置 > 编辑器实时预览功能) |
| **ATK_FRONTEND_REQTIMEOUT** | `15000` | 请求超时 (单位：毫秒) | frontend.reqTimeout (界面配置 > 请求超时) |
| **ATK_FRONTEND_SENDBTN** | `""` | 发送按钮文字 | frontend.sendBtn (界面配置 > 发送按钮文字) |
| **ATK_FRONTEND_UABADGE** | `false` | 用户 UA 徽标 | frontend.uaBadge (界面配置 > 用户 UA 徽标) |
| **ATK_FRONTEND_VERSIONCHECK** | `true` | 版本检测 | frontend.versionCheck (界面配置 > 版本检测) |
| **ATK_FRONTEND_VOTE** | `true` | 投票按钮 | frontend.vote (界面配置 > 投票按钮) |
| **ATK_FRONTEND_VOTEDOWN** | `false` | 反对按钮 | frontend.voteDown (界面配置 > 反对按钮) |


## 服务器

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_HTTP_BODY_LIMIT** | `100` | 请求体大小限制 (单位：MB) | http.body_limit (服务器 > 请求体大小限制) |
| **ATK_HTTP_PROXY_HEADER** | `""` | 代理标头名 (当使用 CDN 时填写 `X-Forwarded-For` 获取用户真实 IP) | http.proxy_header (服务器 > 代理标头名) |


## 图片上传

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_IMG_UPLOAD_ENABLED** | `true` | 启用图片上传 | img_upload.enabled (图片上传 > 启用图片上传) |
| **ATK_IMG_UPLOAD_MAX_SIZE** | `5` | 图片大小限制 (单位：MB) | img_upload.max_size (图片上传 > 图片大小限制) |
| **ATK_IMG_UPLOAD_PATH** | `"./data/artalk-img/"` | 图片存放路径 | img_upload.path (图片上传 > 图片存放路径) |
| **ATK_IMG_UPLOAD_PUBLIC_PATH** | `<nil>` | 图片链接基础路径 (默认为 "/static/images/") | img_upload.public_path (图片上传 > 图片链接基础路径) |
| **ATK_IMG_UPLOAD_UPGIT_DEL_LOCAL** | `true` | 上传后删除本地的图片 | img_upload.upgit.del_local (图片上传 > Upgit 配置 > 上传后删除本地的图片) |
| **ATK_IMG_UPLOAD_UPGIT_ENABLED** | `false` | 启用 Upgit | img_upload.upgit.enabled (图片上传 > Upgit 配置 > 启用 Upgit) |
| **ATK_IMG_UPLOAD_UPGIT_EXEC** | `"upgit -c <upgit配置文件路径> -t /artalk-img"` | 命令行参数 | img_upload.upgit.exec (图片上传 > Upgit 配置 > 命令行参数) |


## IP 属地

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_IP_REGION_DB_PATH** | `"./data/ip2region.xdb"` | 数据文件路径 (.xdb 格式) | ip_region.db_path (IP 属地 > 数据文件路径) |
| **ATK_IP_REGION_ENABLED** | `false` | 启用 IP 属地展示 | ip_region.enabled (IP 属地 > 启用 IP 属地展示) |
| **ATK_IP_REGION_PRECISION** | `"province"` | 显示精度 (可选：`["province", "city", "country"]`) | ip_region.precision (IP 属地 > 显示精度) |


## 日志

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_LOG_ENABLED** | `true` | 启用日志 | log.enabled (日志 > 启用日志) |
| **ATK_LOG_FILENAME** | `"./data/artalk.log"` | 日志文件路径 | log.filename (日志 > 日志文件路径) |


## 评论审核

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_MODERATOR_AKISMET_KEY** | `""` | Akismet Key (Akismet 反垃圾服务，https://akismet.com) | moderator.akismet_key (评论审核 > Akismet Key) |
| **ATK_MODERATOR_ALIYUN_ACCESS_KEY_ID** | `""` | AccessKeyId | moderator.aliyun.access_key_id (评论审核 > 阿里云内容安全 > AccessKeyId) |
| **ATK_MODERATOR_ALIYUN_ACCESS_KEY_SECRET** | `""` | AccessKeySecret | moderator.aliyun.access_key_secret (评论审核 > 阿里云内容安全 > AccessKeySecret) |
| **ATK_MODERATOR_ALIYUN_ENABLED** | `false` | 启用 | moderator.aliyun.enabled (评论审核 > 阿里云内容安全 > Enabled) |
| **ATK_MODERATOR_ALIYUN_REGION** | `"cn-shanghai"` | Region | moderator.aliyun.region (评论审核 > 阿里云内容安全 > Region) |
| **ATK_MODERATOR_API_FAIL_BLOCK** | `false` | API 请求错误时拦截 (关闭此项当请求错误时让评论放行) | moderator.api_fail_block (评论审核 > API 请求错误时拦截) |
| **ATK_MODERATOR_KEYWORDS_ENABLED** | `false` | 启用 | moderator.keywords.enabled (评论审核 > 关键词过滤 > Enabled) |
| **ATK_MODERATOR_KEYWORDS_FILE_SEP** | `"\n"` | 词库文件内容分割符 (例如填写 "\n" 文件中一行一个关键词) | moderator.keywords.file_sep (评论审核 > 关键词过滤 > 词库文件内容分割符) |
| **ATK_MODERATOR_KEYWORDS_FILES** | `[./data/词库_1.txt]` | 词库文件 (支持多个词库文件) | moderator.keywords.files (评论审核 > 关键词过滤 > 词库文件) |
| **ATK_MODERATOR_KEYWORDS_PENDING** | `false` | 匹配成功设为待审状态 | moderator.keywords.pending (评论审核 > 关键词过滤 > 匹配成功设为待审状态) |
| **ATK_MODERATOR_KEYWORDS_REPLACE_TO** | `"x"` | 替换字符 | moderator.keywords.replace_to (评论审核 > 关键词过滤 > 替换字符) |
| **ATK_MODERATOR_PENDING_DEFAULT** | `false` | 默认待审 (发表新评论需要后台人工审核后才能显示) | moderator.pending_default (评论审核 > 默认待审) |
| **ATK_MODERATOR_TENCENT_ENABLED** | `false` | 启用 | moderator.tencent.enabled (评论审核 > 腾讯云文本内容安全 > Enabled) |
| **ATK_MODERATOR_TENCENT_REGION** | `"ap-guangzhou"` | Region | moderator.tencent.region (评论审核 > 腾讯云文本内容安全 > Region) |
| **ATK_MODERATOR_TENCENT_SECRET_ID** | `""` | SecretId | moderator.tencent.secret_id (评论审核 > 腾讯云文本内容安全 > SecretId) |
| **ATK_MODERATOR_TENCENT_SECRET_KEY** | `""` | SecretKey | moderator.tencent.secret_key (评论审核 > 腾讯云文本内容安全 > SecretKey) |


## SSL

| 环境变量 | 默认值 | 描述 | 路径 |
| --- | --- | --- | --- |
| **ATK_SSL_CERT_PATH** | `""` | 证书文件路径 | ssl.cert_path (SSL > 证书文件路径) |
| **ATK_SSL_ENABLED** | `false` | 启用 SSL | ssl.enabled (SSL > 启用 SSL) |
| **ATK_SSL_KEY_PATH** | `""` | 密钥文件路径 | ssl.key_path (SSL > 密钥文件路径) |

<!-- /env-variables -->
</div>

## 开发者

通过 `ATK_DEBUG=1` 启用开发者模式。

本文档由 Artalk 开发者工具自动生成，执行 `make update-conf-docs` 根据模版配置文件 ([conf/artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/tree/master/conf/artalk.example.zh-CN.yml)) 中的注释更新本页内容。

相关代码位于 [@ArtalkJS/Artalk:/internal/config/](https://github.com/ArtalkJS/Artalk/tree/master/internal/config/) 目录。
