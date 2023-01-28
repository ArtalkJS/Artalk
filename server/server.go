package server

import (
	"net/http"

	_ "github.com/ArtalkJS/Artalk/docs/swagger"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/ArtalkJS/Artalk/server/common"
	h "github.com/ArtalkJS/Artalk/server/handler"
	"github.com/ArtalkJS/Artalk/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

// @Title          Artalk API
// @Version        1.0
// @Description    This is an Artalk server.
// @BasePath  	   /api/

// @Contact.name   API Support
// @Contact.url    https://artalk.js.org
// @Contact.email  artalkjs@gmail.com

// @License.name   LGPL-3.0
// @License.url    https://github.com/ArtalkJS/Artalk/blob/master/LICENSE

// @SecurityDefinitions.apikey ApiKeyAuth
// @In   header
// @Name Authorization
// @Description  "Type 'Bearer TOKEN' to correctly set the API Key"
func Init(app *fiber.App) {
	swaggerDocs(app)

	cors(app)
	actionLimit(app)

	if config.Instance.Debug {
		app.Use(pprof.New())
	}

	api := app.Group("/api", middleware.SiteOriginMiddleware())
	{
		h.CommentAdd(api)
		h.CommentGet(api)
		h.Vote(api)
		h.PV(api)
		h.Stat(api)
		h.MarkRead(api)
		h.ImgUpload(api)

		h.Conf(api)
		h.Version(api)

		// captcha
		h.Captcha(api)

		// user
		h.UserGet(api)
		h.UserLogin(api)
		h.UserLoginStatus(api)
		h.UserLogout(api)

		// admin
		admin(api)
	}

	index(app)

	static(app)
	uploadedStatic(app)
}

func admin(f fiber.Router) {
	admin := f.Group("/admin", middleware.AdminOnlyMiddleware())
	{
		h.AdminCommentEdit(admin)
		h.AdminCommentDel(admin)

		h.AdminPageGet(admin)
		h.AdminPageEdit(admin)
		h.AdminPageDel(admin)
		h.AdminPageFetch(admin)

		h.AdminSiteGet(admin)
		h.AdminSiteAdd(admin)
		h.AdminSiteEdit(admin)
		h.AdminSiteDel(admin)

		h.AdminUserGet(admin)
		h.AdminUserAdd(admin)
		h.AdminUserEdit(admin)
		h.AdminUserDel(admin)

		h.AdminCacheWarm(admin)
		h.AdminCacheFlush(admin)

		h.AdminSendMail(admin)
		h.AdminVoteSync(admin)

		h.AdminSettingGet(admin)
		h.AdminSettingSave(admin)
		h.AdminSettingTpl(admin)

		h.AdminTransfer(admin)
	}
}

func cors(f fiber.Router) {
	f.Use(middleware.CorsMiddleware())
}

func actionLimit(f fiber.Router) {
	f.Use(middleware.ActionLimitMiddleware(middleware.ActionLimitConf{
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

func uploadedStatic(f fiber.Router) {
	if config.Instance.ImgUpload.Path == "" {
		config.Instance.ImgUpload.Path = "./data/artalk-img/"
		logrus.Warn("[Image Upload] img_upload.path is not configured, using the default value: " + config.Instance.ImgUpload.Path)
	}

	// 图片上传静态资源可访问路径
	f.Static(config.IMG_UPLOAD_PUBLIC_PATH, config.Instance.ImgUpload.Path)
}

func swaggerDocs(f fiber.Router) {
	f.Get("/swagger/*", swagger.HandlerDefault)
}
