package authcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"sso/app/controllers/api"
	"sso/app/middlewares/jwt"
	"sso/app/models"
	"sso/config/env"
	"sso/repositories/user_repository"
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

type LoginFormVal struct {
	RedirectUrl string
	Errors      []string
}

func New(env *env.Env) *authController {
	return &authController{env: env, AllRepo: api.NewAllRepo(env)}
}

func (auth *authController) Login(ctx *gin.Context) {
	var (
		loginForm LoginForm
		user      *models.User
		err       error
		token     string
	)

	if err = ctx.ShouldBind(&loginForm); err != nil {
		exception.ValidateException(ctx, err, auth.env)

		return
	}

	if user, err = auth.UserRepo.FindByEmail(loginForm.UserName); err != nil {
		log.Debug().Err(err).Msg("auth.UserRepo.FindByEmail")
	}

	if user == nil {
		exception.Unauthorized(ctx)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		exception.Unauthorized(ctx)
		return
	}

	auth.UserRepo.UpdateLastLoginAt(user)
	pretty, _ := auth.UserRepo.LoadUserRoleAndPermissionPretty(user, "sso")
	user.CurrentRoles = pretty.Roles
	user.CurrentPermissions = pretty.Permissions
	auth.env.Auth().SetUser(user)

	if !auth.env.Auth().HasRole("sso") {
		exception.Forbidden(ctx)
		return
	}

	token, err = jwt.GenerateToken(user, auth.env)
	if err != nil {
		log.Error().Err(err).Msg("jwt.GenerateToken")
		exception.InternalErrorWithMsg(ctx, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "data": gin.H{
		"token":    token,
		"lifetime": auth.env.Config().JwtExpiresSeconds,
	}})
}

func (auth *authController) Logout(c *gin.Context) {
	jwt.AddToBlacklist(auth.env.Config().JwtExpiresSeconds, jwt.GetBearerToken(c), auth.env)

	c.JSON(204, nil)
}

func (auth *authController) Info(c *gin.Context) {
	var (
		user   = c.MustGet("user").(*models.User)
		err    error
		pretty *user_repository.UserWithRBAC
	)

	if pretty, err = auth.UserRepo.LoadUserRoleAndPermissionPretty(user, "sso"); err != nil {
		log.Fatal().Err(err).Msg("authController.Info")
		exception.InternalError(c)
		return
	}

	c.JSON(200, gin.H{"data": pretty})
}
