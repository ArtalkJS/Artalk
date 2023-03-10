package captcha

import (
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io"

	"github.com/ArtalkJS/Artalk/internal/config"
)

type Map = map[string]interface{}

type CaptchaPayload struct {
	CheckValue string
	UserID     string
	UserIP     string
}

type Captcha interface {
	Check(CaptchaPayload) (bool, error)
	PageParams() Map
}

//go:embed pages/*
var pages embed.FS

func GetIFrameHTML(t config.CaptchaType) ([]byte, error) {
	f, err := pages.Open("pages/" + string(t) + ".html")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func NewCaptcha(t config.CaptchaType) Captcha {
	switch t {
	case config.TypeTurnstile:
		return NewTurnstileCaptcha(&config.Instance.Captcha.Turnstile)
	case config.TypeReCaptcha:
		return NewReCaptcha(&config.Instance.Captcha.ReCaptcha)
	case config.TypeHCaptcha:
		return NewHCaptcha(&config.Instance.Captcha.HCaptcha)
	case config.TypeGeetest:
		return NewGeetestCaptcha(&config.Instance.Captcha.Geetest)
	default:
		panic("Unknown captcha type")
	}
}

// hmac-sha256 加密
func hmacEncode(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
