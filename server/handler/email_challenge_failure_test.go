package handler_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestEmailChallengeStorageFailures(t *testing.T) {
	for _, operation := range []string{"query", "create", "update", "verify query", "verify update"} {
		t.Run(operation, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			app.Conf().Auth.Email.Enabled = true
			app.Conf().Captcha.Enabled = false
			handler.AuthEmailSend(app.App, router)
			handler.AuthEmailLogin(app.App, router)
			storage := app.Dao().DB()
			challenge := entity.UserEmailVerify{Email: "reader@example.org", Code: "654321", ExpiresAt: time.Now().Add(time.Minute), Model: gorm.Model{UpdatedAt: time.Now().Add(-2 * time.Minute)}}
			if operation != "create" {
				require.NoError(t, storage.Create(&challenge).Error)
			}
			reject := func(tx *gorm.DB) {
				if tx.Statement.Table == "atk_user_email_verifies" {
					tx.AddError(errors.New("challenge database unavailable"))
				}
			}
			switch operation {
			case "query", "verify query":
				require.NoError(t, storage.Callback().Query().Before("gorm:query").Register("test:challenge", reject))
			case "create":
				require.NoError(t, storage.Callback().Create().Before("gorm:create").Register("test:challenge", reject))
			default:
				require.NoError(t, storage.Callback().Update().Before("gorm:update").Register("test:challenge", reject))
			}
			endpoint := "/auth/email/send"
			if operation == "verify query" || operation == "verify update" {
				endpoint = "/auth/email/login"
			}
			flowRequest(t, router, "POST", endpoint, map[string]string{"email": challenge.Email, "code": challenge.Code}, "", 500)
			require.NoError(t, storage.Callback().Query().Remove("test:challenge"))
			require.NoError(t, storage.Callback().Create().Remove("test:challenge"))
			require.NoError(t, storage.Callback().Update().Remove("test:challenge"))
			var saved entity.UserEmailVerify
			if operation == "create" {
				require.ErrorIs(t, storage.First(&saved).Error, gorm.ErrRecordNotFound)
			} else {
				require.NoError(t, storage.First(&saved, challenge.ID).Error)
				require.Equal(t, challenge.Code, saved.Code)
				require.Zero(t, saved.Attempts)
			}
		})
	}
}

func TestEmailChallengeTemplates(t *testing.T) {
	for _, template := range []string{"custom", "missing", "directory"} {
		t.Run(template, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			app.Conf().Auth.Email.Enabled = true
			app.Conf().Captcha.Enabled = false
			path := filepath.Join(t.TempDir(), "verification.html")
			if template == "custom" {
				require.NoError(t, os.WriteFile(path, []byte("Code: {{code}}"), 0600))
			}
			if template == "directory" {
				require.NoError(t, os.Mkdir(path, 0700))
			}
			app.Conf().Auth.Email.VerifyTpl = path
			app.Conf().Auth.Email.VerifySubject = "Code {{code}}"
			handler.AuthEmailSend(app.App, router)
			flowRequest(t, router, "POST", "/auth/email/send", map[string]string{"email": "invalid"}, "", 400)
			flowRequest(t, router, "POST", "/auth/email/send", map[string]string{"email": "reader@example.org"}, "", 200)
			var code entity.UserEmailVerify
			require.NoError(t, app.Dao().DB().First(&code).Error)
			require.Len(t, code.Code, 6)
			app.Conf().Auth.Email.Enabled = false
			flowRequest(t, router, "POST", "/auth/email/send", map[string]string{"email": "reader@example.org"}, "", 400)
		})
	}
}
