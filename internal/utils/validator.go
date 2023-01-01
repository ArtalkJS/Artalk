package utils

import "github.com/asaskevich/govalidator"

func ValidateEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func ValidateURL(url string) bool {
	return govalidator.IsRequestURL(url)
}
