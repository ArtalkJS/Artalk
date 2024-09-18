package captcha

import (
	"github.com/artalkjs/artalk/v2/internal/captcha/image_captcha"
)

var _ Checker = (*ImageChecker)(nil)

type ImageChecker struct {
	User *User
}

func NewImageChecker(user *User) *ImageChecker {
	return &ImageChecker{
		User: user,
	}
}

func (c *ImageChecker) Type() CaptchaType {
	return Image
}

func (c *ImageChecker) Check(value string) (bool, error) {
	return image_captcha.CheckImageCaptchaCode(c.User.IP, value), nil
}

func (c *ImageChecker) Get() ([]byte, error) {
	return image_captcha.GetNewImageCaptchaBase64(c.User.IP)
}
