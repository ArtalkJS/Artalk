package handler

import (
	"github.com/artalkjs/artalk/v2/internal/auth"
	"github.com/artalkjs/artalk/v2/internal/auth/gothic_fiber"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

func SocialLoginGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Enabled {
			return common.RespError(c, 404, "Auth Api disabled")
		}

		return handler(c)
	}
}

type ResponseConfAuthProviders struct {
	Providers []auth.AuthProviderInfo `json:"providers" validator:"required"`
	Anonymous bool                    `json:"anonymous" validator:"required"`
}

// @Id            GetSocialLoginProviders
// @Summary       Get Social Login Providers
// @Description   Get social login providers
// @Tags          System
// @Produce       json
// @Success       200  {object}  ResponseConfAuthProviders
// @Success       404  {object}  Map{msg=string}
// @Router        /conf/auth/providers  [get]
func AuthSocialLogin(app *core.App, router fiber.Router) {
	// Load providers
	var providers []goth.Provider

	loadProviders := func() {
		providers = auth.GetProviders(app.Conf())
		goth.ClearProviders()
		goth.UseProviders(providers...)
	}

	loadProviders()
	app.OnConfUpdated().Add(func(e *core.ConfUpdatedEvent) error {
		loadProviders()
		return nil
	})

	// Endpoints
	router.Get("/conf/auth/providers", SocialLoginGuard(app, func(c *fiber.Ctx) error {
		return common.RespData(c, ResponseConfAuthProviders{
			Providers: auth.GetProviderInfo(app.Conf(), providers),
			Anonymous: app.Conf().Auth.Anonymous,
		})
	}))

	router.Get("/auth/:provider", SocialLoginGuard(app, func(c *fiber.Ctx) error {
		return gothic_fiber.BeginAuthHandler(c)
	}))

	router.Get("/auth/:provider/callback", SocialLoginGuard(app, func(c *fiber.Ctx) error {
		provider, err := gothic_fiber.GetProviderName(c)
		if err != nil {
			log.Error("[SocialLogin] ", err)
			return common.RespError(c, 500, "Field to get provider name")
		}

		// Fetch user
		gothUser, err := gothic_fiber.CompleteUserAuth(c)
		if err != nil {
			log.Error("[SocialLogin] ", err)
			return common.RespError(c, 500, "Field to complete user auth")
		}

		// Convert to social user
		socialUser := auth.GetSocialUser(gothUser)
		log.Debug("[SocialLogin] ", socialUser)

		// Find auth identity
		authIdentity := app.Dao().FindAuthIdentityByRemoteUID(provider, socialUser.RemoteUID)

		// No auth identity record, register user
		if authIdentity.IsEmpty() {
			authIdentity, err = auth.RegisterSocialUser(app.Dao(), socialUser)
			if err != nil {
				log.Error("[SocialLogin] ", err)
				return common.RespError(c, 500, "Failed to register user")
			}
		}
		if authIdentity.UserID == 0 {
			return common.RespError(c, 500, "Auth Identity user_id invalid")
		}

		// Find user perform login
		user := app.Dao().FindUserByID(authIdentity.UserID)
		if user.IsEmpty() {
			return common.RespError(c, 500, "Failed to find user")
		}

		// Get user token
		jwtToken, err := common.LoginGetUserToken(user, app.Conf().AppKey, app.Conf().LoginTimeout)
		if err != nil {
			return common.RespError(c, 500, err.Error())
		}

		// Render response
		return auth.ResponseCallbackPage(c, jwtToken)
	}))
}
