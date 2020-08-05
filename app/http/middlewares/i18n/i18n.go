package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
	"sso/config/env"
)

var UserPreferLangKey = "UserPreferLangKey"

func I18nMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var locate = c.GetHeader("I18n-Language")
		if locate == "" {
			locate = "en"
		}

		switch locate {
		case "zh":
			fallthrough
		case "zh-CN":
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				if trans, found := env.GetUniversalTranslator().GetTranslator("zh"); found {
					zh.RegisterDefaultTranslations(v, trans)
					c.Set(UserPreferLangKey, locate)
				}
			}
		default:
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				if trans, found := env.GetUniversalTranslator().GetTranslator("en"); found {
					en.RegisterDefaultTranslations(v, trans)
					c.Set(UserPreferLangKey, locate)
				}
			}
		}
		c.Next()
	}
}
