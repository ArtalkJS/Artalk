package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAuthSSOExchange struct {
	Token string `json:"token" validate:"required"` // External IdP access token (e.g. an Auth0 access token)
}

// ssoUserinfo is the subset of OIDC /userinfo we care about.
// Compatible with Auth0, Keycloak, and any spec-conforming OIDC provider.
type ssoUserinfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Nickname      string `json:"nickname"`
	GivenName     string `json:"given_name"`
}

// @Id            AuthSSOExchange
// @Summary       Exchange an external IdP access token for an Artalk JWT
// @Description   Validates a third-party OIDC access token (currently Auth0)
//                by calling the issuer's /userinfo endpoint, then mints an
//                Artalk session JWT for the user identified by the email
//                claim. Use when the surrounding application already runs
//                OIDC and you want Artalk comments to inherit that session
//                without showing Artalk's own login UI.
// @Tags          Auth
// @Param         body  body  ParamsAuthSSOExchange  true  "External SSO token"
// @Accept        json
// @Produce       json
// @Success       200  {object}  ResponseUserLogin
// @Failure       400  {object}  Map{msg=string}
// @Failure       401  {object}  Map{msg=string}
// @Failure       404  {object}  Map{msg=string}
// @Failure       500  {object}  Map{msg=string}
// @Failure       503  {object}  Map{msg=string}
// @Router        /sso/exchange  [post]
func AuthSSOExchange(app *core.App, router fiber.Router) {
	router.Post("/sso/exchange", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Enabled || !app.Conf().Auth.SSO.Enabled {
			return common.RespError(c, 404, "SSO token exchange disabled")
		}

		var p ParamsAuthSSOExchange
		if ok, resp := common.ParamsDecode(c, &p); !ok {
			return resp
		}

		issuer := strings.TrimSpace(app.Conf().Auth.SSO.Issuer)
		if issuer == "" {
			return common.RespError(c, 500, "SSO issuer not configured")
		}
		issuer = strings.TrimSuffix(issuer, "/")
		if !strings.HasPrefix(issuer, "http://") && !strings.HasPrefix(issuer, "https://") {
			issuer = "https://" + issuer
		}

		// Validate by calling the IdP's /userinfo (OIDC standard). For Auth0,
		// this also enforces token-not-revoked at the IdP. We trust the IdP
		// to verify the signature on its own token; that's the contract of
		// /userinfo. No JWKS handling here keeps the dep tree small.
		req, err := http.NewRequestWithContext(c.Context(), "GET", issuer+"/userinfo", nil)
		if err != nil {
			return common.RespError(c, 500, err.Error())
		}
		req.Header.Set("Authorization", "Bearer "+p.Token)
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return common.RespError(c, 503, "Cannot reach SSO issuer: "+err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return common.RespError(c, 401, "SSO token rejected by issuer")
		}

		var ui ssoUserinfo
		if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
			return common.RespError(c, 502, "Bad /userinfo response: "+err.Error())
		}
		email := strings.ToLower(strings.TrimSpace(ui.Email))
		if email == "" {
			return common.RespError(c, 400, "No email claim on SSO token")
		}

		name := strings.TrimSpace(ui.Nickname)
		if name == "" {
			name = strings.TrimSpace(ui.GivenName)
		}
		if name == "" {
			name = strings.TrimSpace(ui.Name)
		}
		if name == "" {
			name = email
		}

		user, err := app.Dao().FindCreateUser(name, email, "")
		if err != nil {
			return common.RespError(c, 500, "User lookup failed: "+err.Error())
		}

		jwtToken, err := common.LoginGetUserToken(user, app.Conf().AppKey, app.Conf().LoginTimeout)
		if err != nil {
			return common.RespError(c, 500, err.Error())
		}

		return common.RespData(c, ResponseUserLogin{
			Token: jwtToken,
			User:  app.Dao().CookUser(&user),
		})
	}))
}
