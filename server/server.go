package server

import (
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/internal/config"
	"github.com/ArtalkJS/ArtalkGo/internal/pkged"
	"github.com/ArtalkJS/ArtalkGo/server/common"
	h "github.com/ArtalkJS/ArtalkGo/server/handler"
	"github.com/ArtalkJS/ArtalkGo/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/sirupsen/logrus"
)

func Init(app *fiber.App) {
	cors(app)
	actionLimit(app)

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
		PathPrefix: "frontend",
		Browse:     false,
	}))
}

func index(f fiber.Router) {
	f.All("/", func(c *fiber.Ctx) error {
		if _, err := pkged.FS().Open("frontend/sidebar/index.html"); err == nil {
			return c.Redirect("./sidebar/", fiber.StatusFound)
		}

		return c.Status(fiber.StatusOK).JSON(common.GetApiVersionDataMap())
	})
}

func uploadedStatic(f fiber.Router) {
	if config.Instance.ImgUpload.Path == "" {
		config.Instance.ImgUpload.Path = "./data/artalk-img/"
		logrus.Warn("图片上传功能 img_upload.path 未配置，使用默认值：" + config.Instance.ImgUpload.Path)
	}

	// 图片上传静态资源可访问路径
	f.Static(config.IMG_UPLOAD_PUBLIC_PATH, config.Instance.ImgUpload.Path)
}
