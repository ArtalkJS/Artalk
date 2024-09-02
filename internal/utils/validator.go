package utils

import (
	"net/url"

	"github.com/asaskevich/govalidator"
)

func ValidateEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func ValidateURL(v string) bool {
	u, err := url.ParseRequestURI(v)
	if err != nil {
		return false
	}
	if len(u.Scheme) == 0 {
		return false
	}
	// Must https or http
	if u.Scheme != "https" && u.Scheme != "http" {
		return false
	}
	return true
}
