package handler

import (
	"cmp"
	"errors"
	"hash/maphash"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	var locks [64]sync.Mutex
	seed := maphash.MakeSeed()

	router.Post("/auth/email/send", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Email.Enabled {
			return common.RespError(c, 400, "Email auth is not enabled")
		}

		var p RequestAuthEmailSend
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Check email
		p.Email = strings.ToLower(strings.TrimSpace(p.Email))
		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, "Invalid email")
		}

		// Mutex for each email to avoid frequency check fail concurrently
		mutex := &locks[maphash.String(seed, p.Email)%uint64(len(locks))]
		mutex.Lock()
		defer mutex.Unlock()

		// Keep the issuance cooldown after successful consumption. Exhausted
		// challenges cannot be replaced until their original expiry, even from another IP.
		now := time.Now()
		var recent int64
		err := app.Dao().DB().Model(&entity.UserEmailVerify{}).Where(
			"((LOWER(email) = ? OR ip = ?) AND updated_at > ?) OR (LOWER(email) = ? AND attempts >= ? AND code <> '' AND expires_at > ?)",
			p.Email, c.IP(), now.Add(-time.Minute), p.Email, emailVerifyMaxAttempts, now,
		).Count(&recent).Error
		if err != nil {
			return common.RespError(c, 500, "Failed to check email verify cooldown")
		}
		if recent > 0 {
			return common.RespError(c, 400, "Send email verify code too frequently")
		}

		code := utils.RandomStringWithAlphabet(6, "0123456789")
		var emailVerify entity.UserEmailVerify
		err = app.Dao().DB().Where("LOWER(email) = ?", p.Email).First(&emailVerify).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			emailVerify = entity.UserEmailVerify{
				Email: p.Email, Code: code, ExpiresAt: now.Add(5 * time.Minute),
				IP: c.IP(), UA: string(c.Request().Header.UserAgent()),
			}
			err = app.Dao().DB().Create(&emailVerify).Error
		case err == nil:
			// Recheck in the write: a concurrent fifth failure must not be
			// overwritten by a resend that passed the earlier cooldown check.
			// A map explicitly resets Attempts to zero on a new challenge.
			result := app.Dao().DB().Model(&entity.UserEmailVerify{}).
				Where("id = ? AND updated_at <= ? AND (attempts < ? OR code = '' OR expires_at <= ?)",
					emailVerify.ID, now.Add(-time.Minute), emailVerifyMaxAttempts, now).
				Updates(map[string]any{
					"email": p.Email, "code": code, "attempts": 0,
					"expires_at": now.Add(5 * time.Minute), "ip": c.IP(),
					"ua": string(c.Request().Header.UserAgent()),
				})
			err = result.Error
			if err == nil && result.RowsAffected != 1 {
				return common.RespError(c, 400, "Send email verify code too frequently")
			}
		}
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

		// Do not write authentication codes or rendered email bodies to logs.
		emailService.AsyncSendTo(emailSubject, emailBody, p.Email)

		return common.RespSuccess(c)
	}))
}

const emailVerifyMaxAttempts = 5

func CheckEmailVerifyCode(app *core.App, c *fiber.Ctx, email string, code string) (ok bool, resp error) {
	email = strings.ToLower(strings.TrimSpace(email))
	code = strings.TrimSpace(code)
	if email == "" || code == "" || !utils.ValidateEmail(email) {
		return false, common.RespError(c, 400, "Invalid email or code")
	}

	// Compare and consume the same challenge version in one conditional write.
	// A failed match must spend an attempt before another request can compare
	// against that version, including simultaneous requests with a correct code.
	for retry := 0; retry < emailVerifyMaxAttempts; retry++ {
		var challenge entity.UserEmailVerify
		err := app.Dao().DB().Where(
			"LOWER(email) = ? AND expires_at > ? AND attempts < ? AND code <> ''",
			email, time.Now(), emailVerifyMaxAttempts,
		).First(&challenge).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		if err != nil {
			return false, common.RespError(c, 500, "Failed to verify email code")
		}
		matched := challenge.Code == code
		values := map[string]any{"attempts": challenge.Attempts + 1}
		if matched {
			values["code"] = ""
			values["attempts"] = emailVerifyMaxAttempts
		}
		// updated_at is issuance time and is deliberately not changed here.
		result := app.Dao().DB().Model(&entity.UserEmailVerify{}).Where(
			"id = ? AND code = ? AND attempts = ? AND updated_at = ? AND expires_at > ?",
			challenge.ID, challenge.Code, challenge.Attempts, challenge.UpdatedAt, time.Now(),
		).UpdateColumns(values)
		if result.Error != nil {
			return false, common.RespError(c, 500, "Failed to verify email code")
		}
		if result.RowsAffected == 1 {
			if matched {
				return true, nil
			}
			break
		}
		// A concurrent attempt or resend won; retry against the new version.
	}

	return false, common.RespError(c, 400, "Invalid or expired email verify code")
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
