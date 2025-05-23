package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式判断是否合法
	matched, err := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	if !matched || err != nil {
		return false
	}
	return true
}
