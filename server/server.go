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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
)

// @Title          Artalk API
// @Version        2.0
// @Description    Artalk is a modern comment system based on Golang.
// @BasePath       /api/v2

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

		ErrorHandler: common.ErrorHandler,
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

	api := fb.Group("/api/v2")
	{
		h.CommentCreate(app, api)
		h.CommentList(app, api)
		h.CommentGet(app, api)
		h.Vote(app, api)
		h.PagePV(app, api)
		h.Stat(app, api)
		h.NotifyList(app, api)
		h.NotifyReadAll(app, api)
		h.NotifyRead(app, api)
		h.Upload(app, api)

		h.Conf(app, api)
		h.Version(app, api)

		// captcha
		h.Captcha(app, api)

		// user
		h.UserInfo(app, api)
		h.UserLogin(app, api)
		h.UserStatus(app, api)

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

func admin(app *core.App, api fiber.Router) {
	h.CommentUpdate(app, api)
	h.CommentDelete(app, api)
	h.PageList(app, api)
	h.PageUpdate(app, api)
	h.PageDelete(app, api)
	h.PageFetch(app, api)
	h.PageFetchAll(app, api)
	h.PageFetchStatus(app, api)
	h.SiteList(app, api)
	h.SiteCreate(app, api)
	h.SiteUpdate(app, api)
	h.SiteDelete(app, api)
	h.UserList(app, api)
	h.UserCreate(app, api)
	h.UserUpdate(app, api)
	h.UserDelete(app, api)
	h.CacheWarmUp(app, api)
	h.CacheFlush(app, api)
	h.EmailSend(app, api)
	h.VoteSync(app, api)
	h.SettingGet(app, api)
	h.SettingApply(app, api)
	h.SettingTemplate(app, api)
	h.Transfer(app, api)
}

func cors(app *core.App, f fiber.Router) {
	f.Use(middleware.CorsMiddleware(app))
}

func actionLimit(app *core.App, f fiber.Router) {
	f.Use(limiter.ActionLimitMiddleware(app, limiter.ActionLimitConf{}))
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
