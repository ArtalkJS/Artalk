# 社交登录

Artalk 默认只需填写昵称和邮箱即可发表评论，无需验证邮箱。

但有时候，我们希望用户能够使用社交账号登录，以减少用户填写信息的时间，或者提高用户信息的真实性，可以通过启用社交登录来实现这一目的。

社交登录目前支持以下多种方式：

| 登录方式 | 接入文档 | 登录方式 | 接入文档 | 登录方式 | 接入文档 |
| --- | --- | --- | --- | --- | --- |
| Google | [查看](https://developers.google.com/identity/protocols/oauth2) | Microsoft | [查看](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow) | Apple | [查看](https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/configuring_your_webpage_for_sign_in_with_apple) |
| Facebook | [查看](https://developers.facebook.com/docs/facebook-login/) | Twitter | [查看](https://developer.twitter.com/en/docs/basics/authentication/overview) | Discord | [查看](https://discord.com/developers/docs/topics/oauth2) |
| Slack | [查看](https://api.slack.com/authentication/oauth-v2) | Github | [查看](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/) | Tiktok | [查看](https://developers.tiktok.com/doc/login) |
| Steam | [查看](https://partner.steamgames.com/doc/webapi_overview/auth) | WeChat | [查看](https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html) | Line | [查看](https://developers.line.biz/en/docs/line-login/integrate-line-login/) |
| GitLab | [查看](https://docs.gitlab.com/ee/api/oauth2.html) | Gitea | [查看](https://docs.gitea.io/en-us/oauth2-provider/) | Mastodon | [查看](https://docs.joinmastodon.org/api/authentication/) |
| Patreon | [查看](https://docs.patreon.com/#oauth) | Auth0 | [查看](https://auth0.com/docs/connections/social/) | 邮箱密码 | [查看](#邮箱密码登录) |

开启社交登录功能仅需在 [控制中心](./sidebar.md#设置) 找到「社交登录」启用该功能，然后填写对应的配置信息即可。也可以通过 [配置文件](../backend/config.md) 或 [环境变量](../env.md#社交登录) 进行配置。


## 邮箱密码登录

![邮箱登录](/images/auth/email_login.png)

启用邮箱密码登录后，评论框顶部的昵称邮箱输入框将被隐藏，发送按钮将显示为登录按钮。用户点击登录按钮后，将会弹出一个登录框，用户可以输入邮箱和密码登录，登录成功后即可发表评论。并且，用户发表的评论将展示「邮箱已验证」的标识。

<img src="/images/auth/email_verified.png" width="550" alt="邮箱已验证标识">

用户可以通过邮箱注册账号，Artalk 将向用户邮箱发送一封带有验证码的邮件。验证码有效期为 10 分钟，验证码发送频率限制为 1 分钟一次。

![邮箱注册](/images/auth/email_register.png)

支持自定义验证码邮件模板和邮件标题，可在 Artalk 控制中心的设置页面的社交登录找到「邮箱验证邮件标题」、「邮箱验证邮件模板」选项进行设置。在配置文件中，可以通过 `auth.email.verify_subject` 和 `auth.email.verify_tpl` 进行设置：

```yaml
auth:
  enabled: true
  email:
    enabled: true
    verify_subject: "您的验证码是 - {{code}}"
    verify_tpl: default
```

默认模版如下：

```html
您的验证码是：{{code}}。请使用它来验证您的电子邮件并登录到 Artalk。如果您没有请求此操作，请忽略此消息。
```

![跳过登录](/images/auth/login_skip.png)

启用邮箱密码登录功能后，仍然可跳过邮箱验证：登录弹窗底部显示 “跳过，不验证” 按钮，点击后评论框顶部恢复为显示原有的昵称、邮箱、网址三个输入框。在设置中勾选「允许匿名评论」可以更改。

## 账号合并工具

登录后如果检测到相同的邮箱下有多个不同用户名的账号，将会弹出账号合并工具，用户可以选择保留其中一个用户名，该邮箱下的所有评论等数据合并到保留的账号下。原有的账号将被删除，评论显示的用户名将会变更为保留的用户名。

![账号合并工具](/images/auth/merge_accounts.png)

## 多种登录方式

Artalk 支持同时启用多种登录方式，用户可以选择任意一种方式登录。

![多种登录方式](/images/auth/multi_login.png)

如果只启用了唯一一种登录，例如 GitHub 登录，将直接弹出 GitHub 的授权登录页面。

![GitHub 授权弹窗](/images/auth/github_login.png)

接入 GitHub 登录可参考文档：[关于创建 GitHub 应用](https://docs.github.com/zh/apps/creating-github-apps/about-creating-github-apps/about-creating-github-apps)，得到 Client ID 和 Client Secret 后，填写到 Artalk 控制中心的设置页面的社交登录中的「GitHub」选项中即可。

## SSO 令牌交换

如果你的站点被嵌入在一个已经通过外部 OIDC 身份提供商（Auth0、Keycloak、Okta 等）完成用户认证的应用中，可以让 Artalk 复用该登录态，而不必显示自己的登录界面。外层应用用它已经持有的 IdP 访问令牌换取一个 Artalk 会话令牌，这样已在父应用登录的用户无需额外点击或弹窗即可发表评论。

该功能**默认关闭**，需手动开启。可通过 [配置文件](../backend/config.md) 或 [环境变量](../env.md#社交登录) 启用，并填写你所用提供商的 OIDC issuer：

```yaml
auth:
  enabled: true
  sso:
    enabled: true
    issuer: "https://tenant.auth0.com" # 例如 "tenant.auth0.com" 或 "https://tenant.auth0.com"
```

启用后，外层应用将 IdP 访问令牌提交到交换接口：

```
POST /api/v2/sso/exchange
Content-Type: application/json

{ "token": "<外部 IdP 访问令牌>" }
```

Artalk 通过调用 issuer 的 `/userinfo` 端点校验该令牌（OIDC 标准做法——由 IdP 在服务端验证并签名响应，因此 Artalk 侧无需处理密钥），读取其中的 `email` 声明，查找或创建对应用户，并返回与其它登录接口一致的响应结构：

```json
{
  "token": "<Artalk 会话令牌>",
  "user": { "name": "...", "email": "...", "is_admin": false }
}
```

前端在 `Artalk.init()` 运行前将其写入 `localStorage["ArtalkUser"]`，组件即把该用户视为已完全登录（不弹窗，管理员用户也不再要求输入管理密码）：

```js
const res = await fetch('https://comments.example.com/api/v2/sso/exchange', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ token: idpAccessToken }), // 来自你的 IdP 的访问令牌
})
const { token, user } = await res.json()

// 在组件挂载前预填 Artalk 登录态
localStorage.setItem('ArtalkUser', JSON.stringify({ ...user, token }))

Artalk.init({ el: '#Comments', server: 'https://comments.example.com', site: 'Blog' })
```

当 SSO 未启用时接口返回 `404`，IdP 拒绝令牌时返回 `401`，令牌不含邮箱声明时返回 `400`。

::: tip

校验完全依赖 issuer 的 `/userinfo` 拒绝无效或已吊销的令牌，且用户按 `email` 声明匹配。请仅对你信任的 issuer 启用此功能，并确保该 issuer 返回的是已验证的邮箱。

:::

## 插件开发

Artalk 的社交登录功能是通过独立的插件实现并采用 Solid.js 开发，代码可在 [@ArtalkJS/Artalk:ui/plugin-auth](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-auth) 找到。

在控制中心启用社交登录功能后，将自动在前端加载该插件。
