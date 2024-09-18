package handler

import (
	"cmp"
	"os"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/sync"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type RequestAuthEmailSend struct {
	Email string `json:"email" validate:"required"`
}

// @Id           SendVerifyEmail
// @Summary      Send verify email
// @Description  Send email including verify code to user
// @Tags         Auth
// @Param        data  body  RequestAuthEmailSend  true  "The data"
// @Success      200  {object}  Map{msg=string}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /auth/email/send  [post]
func AuthEmailSend(app *core.App, router fiber.Router) {
	mutexMap := sync.NewKeyMutex[string]()

	router.Post("/auth/email/send", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Email.Enabled {
			return common.RespError(c, 400, "Email auth is not enabled")
		}

		var p RequestAuthEmailSend
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Check email
		p.Email = strings.TrimSpace(p.Email)
		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, "Invalid email")
		}

		// Mutex for each email to avoid frequency check fail concurrently
		mutex := mutexMap.GetLock(p.Email)
		mutex.Lock()
		defer mutex.Unlock()

		// check if had send email verify code in 1 minutes
		if err := app.Dao().DB().Model(&entity.UserEmailVerify{}).Where("email = ? OR ip = ?", p.Email, c.IP()).Where("updated_at > ?", time.Now().Add(-time.Minute*1)).First(&entity.UserEmailVerify{}).Error; err == nil {
			return common.RespError(c, 400, "Send email verify code too frequently")
		}

		// Save email verify code
		code := utils.RandomStringWithAlphabet(6, "0123456789")

		var emailVerify entity.UserEmailVerify
		err := app.Dao().DB().Where(entity.UserEmailVerify{
			Email: p.Email,
		}).Assign(entity.UserEmailVerify{
			Code:      code,
			ExpiresAt: time.Now().Add(time.Minute * 5),
			IP:        c.IP(),
			UA:        string(c.Request().Header.UserAgent()),
		}).FirstOrCreate(&emailVerify).Error
		if err != nil {
			return common.RespError(c, 500, "Failed to save email verify")
		}

		// Send email
		emailService, err := core.AppService[*core.EmailService](app)
		if err != nil {
			return common.RespError(c, 500, "Failed to get email service")
		}

		// Email template
		tplParams := map[string]any{"code": code}
		emailSubject := cmp.Or(utils.RenderMustaches(app.Conf().Auth.Email.VerifySubject, tplParams), i18n.T("Your Code - {{code}}", tplParams))
		emailBody := ""

		if tpl := loadAuthEmailVerifyTemplate(&app.Conf().Auth.Email); tpl != "" {
			emailBody = utils.RenderMustaches(tpl, tplParams)
		} else {
			emailBody = i18n.T("Your code is: {{code}}. Use it to verify your email and sign in Artalk. If you didn't request this, simply ignore this message.", tplParams)
		}

		log.Debug("Send email to: ", p.Email, " subject: ", emailSubject, " body: ", emailBody)
		emailService.AsyncSendTo(emailSubject, emailBody, p.Email)

		return common.RespSuccess(c)
	}))
}

func CheckEmailVerifyCode(app *core.App, c *fiber.Ctx, email string, code string) (ok bool, resp error) {
	email = strings.TrimSpace(email)
	code = strings.TrimSpace(code)

	if email == "" || code == "" || !utils.ValidateEmail(email) {
		return false, common.RespError(c, 400, "Invalid email or code")
	}

	// Check email verify code
	var emailVerify entity.UserEmailVerify
	app.Dao().DB().Where(entity.UserEmailVerify{
		Email: email,
		Code:  code,
	}).First(&emailVerify)

	if emailVerify.ID == 0 {
		return false, common.RespError(c, 400, "Invalid email verify code")
	}
	if emailVerify.ExpiresAt.Before(time.Now()) {
		return false, common.RespError(c, 400, "Email verify code expired")
	}

	// Revoke email verify code
	if err := app.Dao().DB().Unscoped().Delete(&emailVerify).Error; err != nil {
		return false, common.RespError(c, 500, "Failed to revoke email verify code")
	}

	return true, nil
}

func loadAuthEmailVerifyTemplate(conf *config.AuthEmailConf) string {
	if conf.VerifyTpl == "" || conf.VerifyTpl == "default" {
		return ""
	}

	// check if tpl file exists
	if !utils.CheckFileExist(conf.VerifyTpl) {
		log.Error("Email template file not exists: ", conf.VerifyTpl)
		return ""
	}

	// read tpl file
	fs, err := os.ReadFile(conf.VerifyTpl)
	if err != nil {
		log.Error("Failed to read email template file: ", conf.VerifyTpl, " err: ", err.Error())
		return ""
	}

	return string(fs)
}
