package captcha

import (
	"github.com/artalkjs/artalk/v2/internal/config"
)

type Map = map[string]interface{}

type CheckerConf struct {
	config.CaptchaConf
	User User
}

type User struct {
	ID string
	IP string
}

type CaptchaType = int

const (
	Image CaptchaType = iota
	IFrame
)

type Checker interface {
	Type() CaptchaType
	Check(value string) (bool, error)
	Get() ([]byte, error)
}

func NewCaptchaChecker(conf *CheckerConf) Checker {
	switch conf.CaptchaType {
	case config.TypeImage:
		return NewImageChecker(&conf.User)
	case config.TypeTurnstile:
		return NewTurnstileChecker(&conf.Turnstile, &conf.User)
	case config.TypeReCaptcha:
		return NewReCaptchaChecker(&conf.ReCaptcha, &conf.User)
	case config.TypeHCaptcha:
		return NewHCaptchaChecker(&conf.HCaptcha, &conf.User)
	case config.TypeGeetest:
		return NewGeetestChecker(&conf.Geetest, &conf.User)
	default:
		panic("Unknown captcha type")
	}
}
