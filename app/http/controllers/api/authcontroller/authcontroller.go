package authcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
)

type LoginForm struct {
	UserName string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func New(env *env.Env) *authController {
	return &authController{env: env}
}

type authController struct {
	env *env.Env
}

type LoginFormVal struct {
	RedirectUrl string
	Errors      []string
}

func (auth *authController) Login(ctx *gin.Context) {
	var loginForm LoginForm

	if err := ctx.ShouldBind(&loginForm); err != nil {
		exception.ValidateException(ctx, err, auth.env)

		return
	}

	user := models.User{}.FindByEmail(loginForm.UserName, auth.env)
	printErrorBack := func() {
		ctx.JSON(401, gin.H{"code": 401, "msg": "Unauthorized!"})
	}

	if user == nil {
		printErrorBack()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		printErrorBack()
		return
	}

	user.UpdateLastLoginAt(auth.env)

	token := user.GenerateApiToken(auth.env, false)

	ctx.JSON(200, gin.H{"code": 200, "data": gin.H{
		"token":    token,
		"lifetime": auth.env.Config().AccessTokenLifetime,
	}})
}

func (auth *authController) Logout(c *gin.Context) {
	u, exists := c.Get("user")
	if exists {
		user := u.(*models.User)
		// 让用户的api调用不能再使用
		user.GenerateApiToken(auth.env, true)
	}

	c.JSON(204, nil)
}

func (auth *authController) Info(c *gin.Context) {
	userCtx, _ := c.Get("user")
	user := userCtx.(*models.User)
	if err := auth.env.GetDB().Preload("Roles.Permissions").Find(&user).Error; err != nil {
		log.Fatal().Err(err).Msg("authController.Info")
		return
	}
	if !user.TokenExpired(auth.env) {
		c.JSON(200, gin.H{"data": user})
		return
	}

	c.JSON(401, gin.H{"code": 401})
}
