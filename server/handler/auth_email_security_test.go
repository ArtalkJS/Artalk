package handler_test

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestEmailCodeAttemptBudget(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailLogin(app.App, router)
	handler.AuthEmailRegister(app.App, router)
	handler.AuthEmailSend(app.App, router)
	user := app.Dao().FindUserByID(1000)
	challenge := entity.UserEmailVerify{Email: user.Email, Code: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	// Resending must not reset a locked challenge, even after the ordinary one-minute cooldown.
	require.NoError(t, app.Dao().DB().Model(&challenge).UpdateColumn("updated_at", time.Now().Add(-2*time.Minute)).Error)
	for i := 0; i < 5; i++ {
		endpoint := "/auth/email/login"
		if i%2 == 1 {
			endpoint = "/auth/email/register"
		}
		resp, err := securityRequest(router, "POST", endpoint, map[string]string{
			"email": strings.ToUpper(user.Email), "code": "000000", "password": "new-password",
		}, "", fmt.Sprintf("192.0.2.%d", i+1))
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, 400, resp.StatusCode)
	}
	resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": challenge.Code}, "", "198.51.100.1")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode, "a correct code must fail after five guesses across routes and IPs")
	resp, err = securityRequest(router, "POST", "/auth/email/send", map[string]string{"email": strings.ToUpper(user.Email)}, "", "198.51.100.2")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode, "resending must not bypass the exhausted challenge cooldown")
}

func TestEmailCodeNormalLoginAndReplay(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailLogin(app.App, router)
	user := app.Dao().FindUserByID(1000)
	challenge := entity.UserEmailVerify{Email: user.Email, Code: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	for i := 0; i < 4; i++ {
		resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": "000000"}, "", "")
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, 400, resp.StatusCode)
	}
	resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": " " + strings.ToUpper(user.Email) + " ", "code": " 654321 "}, "", "")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
	var result handler.ResponseUserLogin
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Equal(t, user.ID, result.User.ID)
	require.NotEmpty(t, result.Token)
	resp, err = securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": challenge.Code}, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
}

func TestEmailCodeConcurrentConsumption(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	router.Post("/verify", func(c *fiber.Ctx) error {
		if ok, resp := handler.CheckEmailVerifyCode(app.App, c, "concurrent@example.com", "654321"); !ok {
			return resp
		}
		return c.SendStatus(200)
	})
	challenge := entity.UserEmailVerify{Email: "concurrent@example.com", Code: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	var wg sync.WaitGroup
	statuses := make(chan int, 20)
	start := make(chan struct{})
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			resp, err := securityRequest(router, "POST", "/verify", nil, "", "")
			if err != nil {
				statuses <- 0
				return
			}
			resp.Body.Close()
			statuses <- resp.StatusCode
		}()
	}
	close(start)
	wg.Wait()
	close(statuses)
	successes := 0
	for status := range statuses {
		require.Contains(t, []int{200, 400}, status)
		if status == 200 {
			successes++
		}
	}
	require.Equal(t, 1, successes, "only one concurrent request may consume the code")
}

func TestEmailCodeExpiredAndPasswordLogin(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailLogin(app.App, router)
	user := app.Dao().FindUserByID(1001)
	require.NoError(t, user.SetPasswordEncrypt("test-password"))
	require.NoError(t, app.Dao().UpdateUser(&user))
	challenge := entity.UserEmailVerify{Email: user.Email, Code: "654321", ExpiresAt: time.Now().Add(-time.Second)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	for _, code := range []string{"654321", "000000"} {
		resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": code}, "", "")
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, 400, resp.StatusCode)
	}
	resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "password": "test-password"}, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}

func TestEmailCodeResendAfterExpiry(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	app.Conf().Email.Enabled = false // no SMTP or external email in this test
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailSend(app.App, router)
	handler.AuthEmailLogin(app.App, router)
	user := app.Dao().FindUserByID(1001)
	challenge := entity.UserEmailVerify{Email: strings.ToUpper(user.Email), Code: "654321", Attempts: 5, ExpiresAt: time.Now().Add(-time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	require.NoError(t, app.Dao().DB().Model(&challenge).UpdateColumn("updated_at", time.Now().Add(-6*time.Minute)).Error)
	resp, err := securityRequest(router, "POST", "/auth/email/send", map[string]string{"email": user.Email}, "", "192.0.2.1")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, app.Dao().DB().First(&challenge, challenge.ID).Error)
	require.Zero(t, challenge.Attempts)
	require.Len(t, challenge.Code, 6)
	require.True(t, challenge.ExpiresAt.After(time.Now()))
	require.Equal(t, user.Email, challenge.Email)
	resp, err = securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": challenge.Code}, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
	// Consumption must not remove the one-minute issuance limit.
	resp, err = securityRequest(router, "POST", "/auth/email/send", map[string]string{"email": user.Email}, "", "192.0.2.2")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
}

func TestEmailCodeConcurrentFailuresAndProfileUpdate(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.UserInfoUpdate(app.App, router)
	challenge := entity.UserEmailVerify{Email: "new-profile@example.com", Code: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(1001), app.Conf().AppKey, 60)
	require.NoError(t, err)
	var wg sync.WaitGroup
	statuses := make(chan int, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := securityRequest(router, "POST", "/user", map[string]string{"name": "userA", "email": challenge.Email, "code": "000000"}, token, "")
			if err != nil {
				statuses <- 0
				return
			}
			resp.Body.Close()
			statuses <- resp.StatusCode
		}()
	}
	wg.Wait()
	close(statuses)
	for status := range statuses {
		require.Equal(t, 400, status)
	}
	require.NoError(t, app.Dao().DB().First(&challenge, challenge.ID).Error)
	require.Equal(t, 5, challenge.Attempts)
	resp, err := securityRequest(router, "POST", "/user", map[string]string{"name": "userA", "email": challenge.Email, "code": challenge.Code}, token, "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
	require.Equal(t, "user_a@qwqaq.com", app.Dao().FindUserByID(1001).Email)
}

func TestEmailCodeResendCannotOverwriteConcurrentLockout(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	app.Conf().Email.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailSend(app.App, router)
	challenge := entity.UserEmailVerify{Email: "resend@example.com", Code: "654321", Attempts: 4, ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	require.NoError(t, app.Dao().DB().Model(&challenge).UpdateColumn("updated_at", time.Now().Add(-2*time.Minute)).Error)
	// Inject the fifth failure between the cooldown read and resend write.
	const callback = "test:concurrent_email_failure"
	require.NoError(t, app.Dao().DB().Callback().Update().Before("gorm:update").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table != "atk_user_email_verifies" {
			return
		}
		require.NoError(t, app.Dao().DB().Callback().Update().Remove(callback))
		require.NoError(t, app.Dao().DB().Model(&challenge).UpdateColumn("attempts", 5).Error)
	}))
	defer app.Dao().DB().Callback().Update().Remove(callback)
	resp, err := securityRequest(router, "POST", "/auth/email/send", map[string]string{"email": challenge.Email}, "", "192.0.2.1")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
	require.NoError(t, app.Dao().DB().First(&challenge, challenge.ID).Error)
	require.Equal(t, 5, challenge.Attempts)
	require.Equal(t, "654321", challenge.Code)
}

func TestEmailCodeCannotConsumeReplacedChallenge(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailLogin(app.App, router)
	user := app.Dao().FindUserByID(1001)
	challenge := entity.UserEmailVerify{Email: user.Email, Code: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	const callback = "test:replace_challenge"
	require.NoError(t, app.Dao().DB().Callback().Update().Before("gorm:update").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table != "atk_user_email_verifies" {
			return
		}
		require.NoError(t, app.Dao().DB().Callback().Update().Remove(callback))
		require.NoError(t, app.Dao().DB().Model(&challenge).Updates(map[string]any{"code": "123456", "attempts": 0}).Error)
	}))
	defer app.Dao().DB().Callback().Update().Remove(callback)
	resp, err := securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": "654321"}, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
	require.NoError(t, app.Dao().DB().First(&challenge, challenge.ID).Error)
	require.Equal(t, "123456", challenge.Code)
	resp, err = securityRequest(router, "POST", "/auth/email/login", map[string]string{"email": user.Email, "code": "123456"}, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}
