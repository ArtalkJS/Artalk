package image_captcha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageCaptcha(t *testing.T) {
	testIP := "127.0.0.1"

	buf, err := GetNewImageCaptchaBase64(testIP)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, buf)
	}

	realCode, hit := captchaStore.Get(CaptchaCachePrefix + testIP)
	if assert.True(t, hit, "captcha cached code should be hit") {
		assert.NotEmpty(t, realCode, "captcha cached code should not be empty")
	}

	t.Run("CheckCorrect", func(t *testing.T) {
		ok := CheckImageCaptchaCode(testIP, realCode.(string))
		assert.True(t, ok, "code should be correct")
	})

	t.Run("CheckIncorrect", func(t *testing.T) {
		ok := CheckImageCaptchaCode(testIP, "123456")
		assert.False(t, ok, "code should be incorrect")
	})

	t.Run("Invalidate", func(t *testing.T) {
		InvalidateImageCaptcha(testIP)
		ok := CheckImageCaptchaCode(testIP, realCode.(string))
		assert.False(t, ok, "code should be incorrect")
	})

	t.Run("Regenerate and Check", func(t *testing.T) {
		// generate x1
		buf, err := GetNewImageCaptchaBase64(testIP)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, buf)
		}

		// check x1
		realCode, hit := captchaStore.Get(CaptchaCachePrefix + testIP)
		if assert.True(t, hit, "captcha cached code should be hit") {
			assert.NotEmpty(t, realCode, "captcha cached code should not be empty")
		}

		ok := CheckImageCaptchaCode(testIP, realCode.(string))
		assert.True(t, ok, "code should be correct")

		// generate x2
		buf, err = GetNewImageCaptchaBase64(testIP)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, buf)
		}

		ok = CheckImageCaptchaCode(testIP, realCode.(string))
		assert.False(t, ok, "code should be incorrect after regenerate")
	})
}
