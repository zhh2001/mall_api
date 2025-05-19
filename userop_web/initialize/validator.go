package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"mall_api/userop_web/global"
)

func InitTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}

		switch local {
		case "en":
			err := enTranslations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		case "zh":
			err := zhTranslations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		default:
			err := enTranslations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		}
		return
	}
	return
}
