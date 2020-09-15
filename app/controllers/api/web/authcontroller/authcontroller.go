package authcontroller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/config/env"
)

type LoginForm struct {
	UserName string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func New(env *env.Env) *authController {
	return &authController{env: env, AllRepo: api.NewAllRepo(env)}
}

type authController struct {
	env *env.Env
	*api.AllRepo
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
	userCtx, _ := c.Get("user")
	user := userCtx.(*models.User)
	if err := auth.env.GetDB().Preload("Roles.Permissions").First(&user).Error; err != nil {
		log.Fatal().Err(err).Msg("authController.Info")
		c.AbortWithError(500, errors.New("internal error"))

		return
	}

	c.JSON(200, gin.H{"data": user})
}
