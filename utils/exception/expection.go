package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"sso/app/http/middlewares/i18n"
	"sso/config/env"
	"sso/utils/form"
)

func ModelNotFound(ctx *gin.Context, modelName string) {
	ctx.AbortWithStatusJSON(404, gin.H{
		"code": 404,
		"msg":  modelName + " not found!",
	})
}

func ValidateException(ctx *gin.Context, errors interface{}, env *env.Env) {
	switch errors.(type) {
	case validator.ValidationErrors:
		err := errors.(validator.ValidationErrors)
		value, _ := ctx.Get(i18n.UserPreferLangKey)
		trans, _ := env.GetUniversalTranslator().GetTranslator(value.(string))
		ctx.AbortWithStatusJSON(422, gin.H{"code": 422, "error": form.ErrorsToMap(err, trans)})
	case form.ValidateErrors:
		ctx.AbortWithStatusJSON(422, gin.H{"code": 422, "error": form.ErrorsToMap(errors, nil)})
	default:
		panic(fmt.Sprintf("not support %v", errors))
	}
}
