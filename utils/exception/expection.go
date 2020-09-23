package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"sso/app/middlewares/i18n"
	"sso/config/env"
	"sso/utils/form"
)

var Unauthorized = func(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "Unauthorized!"})
}

var Forbidden = func(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "Forbidden!"})
}

var InternalError = func(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": "internal error"})
}

var InternalErrorWithMsg = func(ctx *gin.Context, msg string) {
	if gin.Mode() != gin.ReleaseMode {
		ctx.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": msg})
		return
	}

	log.Error().Msg(msg)
	InternalError(ctx)
}

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
		if err, ok := errors.(error); ok {
			InternalErrorWithMsg(ctx, err.Error())
			return
		}
		panic(fmt.Sprintf("not support %v", errors))
	}
}
