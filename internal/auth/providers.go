package auth

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/patreon"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/tiktok"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/wechat"
	"github.com/samber/lo"
)

func GetProviders(conf *config.Config) []goth.Provider {
	providers := []goth.Provider{}

	u, _ := url.Parse(conf.Auth.Callback)
	origin := u.Scheme + "://" + u.Host

	callbackURL := func(provider string) string {
		log.Debug("[SocialLogin] Callback URL: ", fmt.Sprintf("%s/api/v2/auth/%s/callback", origin, provider))
		return fmt.Sprintf("%s/api/v2/auth/%s/callback", origin, provider)
	}

	// @see https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
	if githubConf := conf.Auth.Github; githubConf.Enabled {
		providers = append(providers, github.New(githubConf.ClientID, githubConf.ClientSecret, callbackURL("github"),
			"read:user", "user:email"))
	}
	// @see https://docs.gitlab.com/ee/integration/oauth_provider.html
	if gitlabConf := conf.Auth.Gitlab; gitlabConf.Enabled {
		providers = append(providers, gitlab.New(gitlabConf.ClientID, gitlabConf.ClientSecret, callbackURL("gitlab"),
			"read_user", "email"))
	}
	// @see https://docs.gitea.io/en-us/oauth2-provider/
	if giteaConf := conf.Auth.Gitea; giteaConf.Enabled {
		providers = append(providers, gitea.New(giteaConf.ClientID, giteaConf.ClientSecret, callbackURL("gitea"),
			"read:user"))
	}
	// @see https://developers.google.com/identity/protocols/oauth2
	if googleConf := conf.Auth.Google; googleConf.Enabled {
		providers = append(providers, google.New(googleConf.ClientID, googleConf.ClientSecret, callbackURL("google")))
	}
	// @see https://docs.joinmastodon.org/spec/oauth/
	if mastodonConf := conf.Auth.Mastodon; mastodonConf.Enabled {
		providers = append(providers, mastodon.New(mastodonConf.ClientID, mastodonConf.ClientSecret, callbackURL("mastodon")))
	}
	// @see https://developer.twitter.com/en/docs/authentication/oauth-2-0
	if twitterConf := conf.Auth.Twitter; twitterConf.Enabled {
		providers = append(providers, twitter.New(twitterConf.ClientID, twitterConf.ClientSecret, callbackURL("twitter")))
	}
	// @see https://developers.facebook.com/docs/facebook-login
	if facebookConf := conf.Auth.Facebook; facebookConf.Enabled {
		providers = append(providers, facebook.New(facebookConf.ClientID, facebookConf.ClientSecret, callbackURL("facebook")))
	}
	// @see https://discord.com/developers/docs/topics/oauth2
	if discordConf := conf.Auth.Discord; discordConf.Enabled {
		providers = append(providers, discord.New(discordConf.ClientID, discordConf.ClientSecret, callbackURL("discord"),
			discord.ScopeIdentify, discord.ScopeEmail))
	}
	// @see https://developer.valvesoftware.com/wiki/Steam_Web_API
	if steamConf := conf.Auth.Steam; steamConf.Enabled {
		providers = append(providers, steam.New(steamConf.ApiKey, callbackURL("steam")))
	}
	// @see https://developer.apple.com/documentation/sign_in_with_apple
	if appleConf := conf.Auth.Apple; appleConf.Enabled {
		providers = append(providers, apple.New(appleConf.ClientID, appleConf.ClientSecret, callbackURL("apple"), nil,
			apple.ScopeEmail, apple.ScopeName))
	}
	// @see https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
	if microsoftConf := conf.Auth.Microsoft; microsoftConf.Enabled {
		providers = append(providers, microsoftonline.New(microsoftConf.ClientID, microsoftConf.ClientSecret, callbackURL("microsoftonline"),
			"openid", "offline_access", "profile", "user.read", "email"))
	}
	// @see https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	if wechatConf := conf.Auth.Wechat; wechatConf.Enabled {
		providers = append(providers, wechat.New(wechatConf.ClientID, wechatConf.ClientSecret, callbackURL("wechat"),
			lo.If(strings.HasPrefix(conf.Locale, "zh"), wechat.WECHAT_LANG_CN).Else(wechat.WECHAT_LANG_EN),
		))
	}
	// @see https://developers.tiktok.com/
	if tiktokConf := conf.Auth.Tiktok; tiktokConf.Enabled {
		providers = append(providers, tiktok.New(tiktokConf.ClientID, tiktokConf.ClientSecret, callbackURL("tiktok")))
	}
	// @see https://api.slack.com/authentication/oauth-v2
	if slackConf := conf.Auth.Slack; slackConf.Enabled {
		providers = append(providers, slack.New(slackConf.ClientID, slackConf.ClientSecret, callbackURL("slack")))
	}
	// @see https://developers.line.biz/en/docs/line-login/integrate-line-login/
	if lineConf := conf.Auth.Line; lineConf.Enabled {
		providers = append(providers, line.New(lineConf.ClientID, lineConf.ClientSecret, callbackURL("line")))
	}
	// @see https://www.patreon.com/portal/registration
	if patreonConf := conf.Auth.Patreon; patreonConf.Enabled {
		providers = append(providers, patreon.New(patreonConf.ClientID, patreonConf.ClientSecret, callbackURL("patreon"),
			patreon.ScopeIdentity, patreon.ScopeIdentityEmail))
	}
	// @see https://auth0.com/docs/api/authentication
	if auth0Conf := conf.Auth.Auth0; auth0Conf.Enabled {
		providers = append(providers, auth0.New(auth0Conf.ClientID, auth0Conf.ClientSecret, callbackURL("auth0"),
			auth0Conf.Domain, "openid", "profile", "email"))
	}

	return providers
}
