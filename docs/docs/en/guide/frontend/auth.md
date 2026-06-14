# Social Login

By default, Artalk allows users to post comments by simply filling in their nickname and email without email verification.

However, sometimes we want users to log in with their social accounts to reduce the time spent filling in information or to increase the authenticity of user information. This can be achieved by enabling social login.

Currently, the following social login methods are supported:

| Login Method | Integration Docs | Login Method | Integration Docs | Login Method | Integration Docs |
| ------------ | ---------------- | ------------ | ---------------- | ------------ | ---------------- |
| Google       | [View](https://developers.google.com/identity/protocols/oauth2) | Microsoft | [View](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow) | Apple | [View](https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/configuring_your_webpage_for_sign_in_with_apple) |
| Facebook     | [View](https://developers.facebook.com/docs/facebook-login/) | Twitter    | [View](https://developer.twitter.com/en/docs/basics/authentication/overview) | Discord | [View](https://discord.com/developers/docs/topics/oauth2) |
| Slack        | [View](https://api.slack.com/authentication/oauth-v2) | GitHub     | [View](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/) | TikTok | [View](https://developers.tiktok.com/doc/login) |
| Steam        | [View](https://partner.steamgames.com/doc/webapi_overview/auth) | WeChat     | [View](https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html) | Line | [View](https://developers.line.biz/en/docs/line-login/integrate-line-login/) |
| GitLab       | [View](https://docs.gitlab.com/ee/api/oauth2.html) | Gitea      | [View](https://docs.gitea.io/en-us/oauth2-provider/) | Mastodon | [View](https://docs.joinmastodon.org/api/authentication/) |
| Patreon      | [View](https://docs.patreon.com/#oauth) | Auth0      | [View](https://auth0.com/docs/connections/social/) | Email & Password | [View](#email-password-login) |

To enable social login, simply find the "Social Login" option in the [Dashboard](./sidebar.md#settings), enable it, and fill in the corresponding configuration information. Alternatively, you can configure it through the [configuration file](../backend/config.md) or [environment variables](../env.md#social-login).

## Email & Password Login

![Email Login](/images/auth/email_login.png)

After enabling email and password login, the nickname and email input fields at the top of the comment box will be hidden, and the send button will change to a login button. When the user clicks the login button, a login box will pop up, allowing the user to log in with their email and password. Once logged in, the user can post comments. Comments posted by the user will display a "Email Verified" badge.

<img src="/images/auth/email_verified.png" width="550" alt="Email Verified Badge">

Users can register an account via email, and Artalk will send a verification code to their email. The verification code is valid for 10 minutes, and the frequency of sending verification codes is limited to once per minute.

![Email Registration](/images/auth/email_register.png)

You can customize the verification email template and subject. In the settings page of the Artalk Dashboard under social login, you can find options for "Email Verification Subject" and "Email Verification Template". In the configuration file, you can set them using `auth.email.verify_subject` and `auth.email.verify_tpl`:

```yaml
auth:
  enabled: true
  email:
    enabled: true
    verify_subject: "Your verification code is - {{code}}"
    verify_tpl: default
```

The default template is as follows:

```html
Your verification code is: {{code}}. Please use it to verify your email and log in to Artalk. If you did not request this, please ignore this message.
```

![Skip Login](/images/auth/login_skip.png)

After enabling email and password login, users can still skip email verification: the login popup will show a "Skip, do not verify" button at the bottom, which, when clicked, restores the original nickname, email, and website input fields at the top of the comment box. You can change this by selecting "Allow Anonymous Comments" in the settings.

## Account Merging Tool

If multiple accounts with different usernames but the same email are detected after logging in, an account merging tool will pop up. Users can choose to keep one username, and all comments and data associated with that email will be merged into the retained account. The original accounts will be deleted, and the displayed username on comments will change to the retained username.

![Account Merging Tool](/images/auth/merge_accounts.png)

## Multiple Login Methods

Artalk supports enabling multiple login methods simultaneously, allowing users to choose any method to log in.

![Multiple Login Methods](/images/auth/multi_login.png)

If only one login method is enabled, such as GitHub login, the GitHub authorization login page will pop up directly.

![GitHub Authorization Popup](/images/auth/github_login.png)

For integrating GitHub login, refer to the documentation: [About Creating GitHub Apps](https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app). After obtaining the Client ID and Client Secret, fill them in the "GitHub" option in the social login settings page of the Artalk Dashboard.

## SSO Token Exchange

If your site is embedded in an application that already authenticates users through an external OIDC identity provider (Auth0, Keycloak, Okta, etc.), you can let Artalk reuse that session instead of showing its own login UI. The surrounding application exchanges the IdP access token it already holds for an Artalk session token, so users who are already signed in to the parent app can comment without an extra click or popup.

This feature is **disabled by default** and is opt-in only. Enable it through the [configuration file](../backend/config.md) or [environment variables](../env.md#social-login), filling in the OIDC issuer of your provider:

```yaml
auth:
  enabled: true
  sso:
    enabled: true
    issuer: "https://tenant.auth0.com" # e.g. "tenant.auth0.com" or "https://tenant.auth0.com"
```

Once enabled, the surrounding application posts the IdP access token to the exchange endpoint:

```
POST /api/v2/sso/exchange
Content-Type: application/json

{ "token": "<external IdP access token>" }
```

Artalk validates the token by calling the issuer's `/userinfo` endpoint (the standard OIDC way — the IdP verifies and signs the response server-side, so no key handling is required on Artalk's side), reads the `email` claim, finds or creates the matching user, and returns the same response shape as the other login endpoints:

```json
{
  "token": "<Artalk session token>",
  "user": { "name": "...", "email": "...", "is_admin": false }
}
```

The frontend writes this into `localStorage["ArtalkUser"]` before `Artalk.init()` runs, and the widget then treats the user as fully logged in (no popup, and no admin password prompt for admin users):

```js
const res = await fetch('https://comments.example.com/api/v2/sso/exchange', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ token: idpAccessToken }), // the access token from your IdP
})
const { token, user } = await res.json()

// Pre-populate the Artalk session before the widget mounts
localStorage.setItem('ArtalkUser', JSON.stringify({ ...user, token }))

Artalk.init({ el: '#Comments', server: 'https://comments.example.com', site: 'Blog' })
```

The endpoint returns `404` when SSO is not enabled, `401` when the IdP rejects the token, and `400` when the token carries no email claim.

::: tip

Validation relies entirely on the issuer's `/userinfo` rejecting invalid or revoked tokens, and the user is matched by the `email` claim. Only enable this for an issuer you trust, and make sure that issuer returns verified emails.

:::

## Plugin Development

The social login feature of Artalk is implemented through an independent plugin developed using Solid.js. The code can be found in [@ArtalkJS/Artalk:ui/plugin-auth](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-auth).

After enabling the social login feature in the Dashboard, the plugin will be automatically loaded on the front end.
