package authcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
)

type LoginForm struct {
	UserName string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type authController struct {
	env *env.Env
	*api.AllRepo
}

func NewAuthController(env *env.Env) *authController {
	return &authController{env: env, AllRepo: api.NewAllRepo(env)}
}

func (auth *authController) Logout(c *gin.Context) {
	u, exists := c.Get("user")
	if exists {
		user := u.(*models.User)
		auth.UserRepo.ForceLogout(user)
	}

	c.JSON(204, nil)
}

func (auth *authController) Info(c *gin.Context) {
	type Uri struct {
		Project string `json:"project" uri:"project"`
	}

	var (
		data interface{}
		err  error
		user = c.MustGet("user").(*models.User)
		uri  = &Uri{}
	)

	if err = c.ShouldBindUri(uri); err != nil {
		log.Error().Err(err).Msg("authController.Info")
		exception.InternalErrorWithMsg(c, err.Error())

		return
	}

	if data, err = auth.UserRepo.LoadUserRoleAndPermissionPretty(user, uri.Project); err != nil {
		log.Error().Err(err).Msg("authController.Info")
		exception.InternalErrorWithMsg(c, err.Error())

		return
	}

	c.JSON(200, gin.H{"data": data})
}
