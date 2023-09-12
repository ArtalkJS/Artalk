package server

import (
	"fmt"
	"io"
	"net/http"

	_ "github.com/ArtalkJS/Artalk/docs/swagger"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/ArtalkJS/Artalk/server/common"
	h "github.com/ArtalkJS/Artalk/server/handler"
	"github.com/ArtalkJS/Artalk/server/middleware"
	"github.com/ArtalkJS/Artalk/server/middleware/limiter"
	"github.com/ArtalkJS/Artalk/server/middleware/site_origin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
)

// @Title          Artalk API
// @Version        1.0
// @Description    This is an Artalk server.
// @BasePath  	   /api/

// @Contact.name   API Support
// @Contact.url    https://artalk.js.org
// @Contact.email  artalkjs@gmail.com

// @License.name   MIT
// @License.url    https://github.com/ArtalkJS/Artalk/blob/master/LICENSE

// @SecurityDefinitions.apikey ApiKeyAuth
// @In   header
// @Name Authorization
// @Description  "Type 'Bearer TOKEN' to correctly set the API Key"
func Serve(app *core.App) (*fiber.App, error) {
	// create fiber app
	// @see https://docs.gofiber.io/api/fiber#config
	fb := fiber.New(fiber.Config{
		// @see https://github.com/gofiber/fiber/issues/426
		// @see https://github.com/gofiber/fiber/issues/185
		Immutable: true,

		// TODO add config to set ProxyHeader
		ProxyHeader:        "X-Forwarded-For",
		EnableIPValidation: true,
	})

	// logger
	fb.Use(fiber_logger.New(fiber_logger.Config{
		Format: "[${status}] ${method} ${path} ${latency} ${ip} ${reqHeader:X-Request-ID} ${referer} ${ua}\n",
		Output: io.Discard,
		Done: func(c *fiber.Ctx, logString []byte) {
			statusOK := c.Response().StatusCode() >= 200 && c.Response().StatusCode() <= 302
			if !statusOK {
				log.StandardLogger().WriterLevel(log.ErrorLevel).Write(logString)
			} else {
				log.StandardLogger().WriterLevel(log.DebugLevel).Write(logString)
			}
		},
	}))

	swaggerDocs(fb)

	cors(app, fb)
	actionLimit(app, fb)

	if app.Conf().Debug {
		fb.Use(pprof.New())
	}

	api := fb.Group("/api", site_origin.SiteOriginMiddleware(app))
	{
		h.CommentAdd(app, api)
		h.CommentGet(app, api)
		h.Vote(app, api)
		h.PV(app, api)
		h.Stat(app, api)
		h.MarkRead(app, api)
		h.ImgUpload(app, api)

		h.Conf(app, api)
		h.Version(app, api)

		// captcha
		h.Captcha(app, api)

		// user
		h.UserGet(app, api)
		h.UserLogin(app, api)
		h.UserLoginStatus(app, api)
		h.UserLogout(app, api)

		// admin
		admin(app, api)
	}

	index(fb)

	static(fb)
	uploadedStatic(app, fb)

	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		_ = fb.Shutdown()
		return nil
	})

	// listen serves HTTP requests
	listenAddr := fmt.Sprintf("%s:%d", app.Conf().Host, app.Conf().Port)
	if app.Conf().SSL.Enabled {
		return fb, fb.ListenTLS(listenAddr, app.Conf().SSL.CertPath, app.Conf().SSL.KeyPath)
	} else {
		return fb, fb.Listen(listenAddr)
	}
}

func admin(app *core.App, f fiber.Router) {
	admin := f.Group("/admin", middleware.AdminOnlyMiddleware(app))
	{
		h.AdminCommentEdit(app, admin)
		h.AdminCommentDel(app, admin)

		h.AdminPageGet(app, admin)
		h.AdminPageEdit(app, admin)
		h.AdminPageDel(app, admin)
		h.AdminPageFetch(app, admin)

		h.AdminSiteGet(app, admin)
		h.AdminSiteAdd(app, admin)
		h.AdminSiteEdit(app, admin)
		h.AdminSiteDel(app, admin)

		h.AdminUserGet(app, admin)
		h.AdminUserAdd(app, admin)
		h.AdminUserEdit(app, admin)
		h.AdminUserDel(app, admin)

		h.AdminCacheWarm(app, admin)
		h.AdminCacheFlush(app, admin)

		h.AdminSendMail(app, admin)
		h.AdminVoteSync(app, admin)

		h.AdminSettingGet(app, admin)
		h.AdminSettingSave(app, admin)
		h.AdminSettingTpl(app, admin)

		h.AdminTransfer(app, admin)
	}
}

func cors(app *core.App, f fiber.Router) {
	f.Use(middleware.CorsMiddleware(app))
}

func actionLimit(app *core.App, f fiber.Router) {
	f.Use(limiter.ActionLimitMiddleware(app, limiter.ActionLimitConf{
		// 启用操作限制路径白名单
		ProtectPaths: []string{
			"/api/add",
			"/api/login",
			"/api/vote",
			"/api/img-upload",
		},
	}))
}

func static(f fiber.Router) {
	f.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(pkged.FS()),
		PathPrefix: "public",
		Browse:     false,
	}))
}

func index(f fiber.Router) {
	f.All("/", func(c *fiber.Ctx) error {
		if _, err := pkged.FS().Open("public/sidebar/index.html"); err == nil {
			return c.Redirect("./sidebar/", fiber.StatusFound)
		}

		return c.Status(fiber.StatusOK).JSON(common.GetApiVersionDataMap())
	})
}

func uploadedStatic(app *core.App, f fiber.Router) {
	// 图片上传静态资源可访问路径
	f.Static(config.IMG_UPLOAD_PUBLIC_PATH, app.Conf().ImgUpload.Path)
}

func swaggerDocs(f fiber.Router) {
	f.Get("/swagger/*", swagger.HandlerDefault)
}
