package authcontroller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
)

type LoginForm struct {
	UserName    string `form:"email" binding:"required"`
	Password    string `form:"password" binding:"required"`
	RedirectUrl string `form:"redirect_url"`
}

func New(env *env.Env) *authController {
	return &authController{env: env, AllRepo: api.NewAllRepo(env)}
}

type authController struct {
	env *env.Env
	*api.AllRepo
}

type LoginFormVal struct {
	RedirectUrl string
	Errors      []string
}

func (*authController) LoginForm(ctx *gin.Context) {
	redirectUrl := ctx.Query("redirect_url")
	ctx.HTML(http.StatusOK, "login.tmpl", LoginFormVal{
		RedirectUrl: redirectUrl,
	})
}

func (auth *authController) Login(ctx *gin.Context) {
	var loginForm LoginForm
	redirectUrl := ctx.Query("redirect_url")
	if err := ctx.ShouldBind(&loginForm); err != nil {
		exception.ValidateException(ctx, err, auth.env)

		return
	}

	user, err := auth.UserRepo.FindByEmail(loginForm.UserName, auth.env)
	if err != nil {
		log.Error().Err(err).Msg("auth.UserRepo.FindByEmail")
	}
	printErrorBack := func() {
		ctx.HTML(200, "login.tmpl", LoginFormVal{
			RedirectUrl: redirectUrl,
			Errors:      []string{"username or password error."},
		})
	}

	if user == nil {
		printErrorBack()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		printErrorBack()
		return
	}

	session := sessions.Default(ctx)
	session.Set("user", user)
	err = session.Save()
	if err != nil {
		log.Debug().Err(err).Msg("authController.Login")
	}

	auth.UserRepo.UpdateLastLoginAt(user)

	if loginForm.RedirectUrl == "" {
		ctx.Redirect(302, "/")
		return
	}

	token := auth.UserRepo.GenerateAccessToken(user)

	ctx.Redirect(302, loginForm.RedirectUrl+"?access_token="+token)
}

func (auth *authController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user, ok := session.Get("user").(*models.User)
	session.Clear()
	session.Save()
	if ok {
		auth.UserRepo.ForceLogout(user)
	}

	c.Redirect(302, "/login")
}

func (auth *authController) SelectSystem(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(*models.User)
	c.HTML(200, "select_system.tmpl", struct {
		AccessToken string
	}{
		AccessToken: auth.UserRepo.GenerateAccessToken(user),
	})
}

func (auth *authController) AccessToken(c *gin.Context) {
	var jsonData struct {
		AccessToken string `json:"access_token"`
	}
	err := c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	log.Debug().Msg(jsonData.AccessToken)

	conn := auth.env.RedisPool().Get()

	defer conn.Close()

	id, err := redis.Int(conn.Do("GET", jsonData.AccessToken))
	log.Debug().Err(err).Interface("id", id).Msg("authController.AccessToken")
	if err == nil {
		user, _ := auth.UserRepo.FindById(uint(id))
		if user != nil {
			do, err := conn.Do("DEL", jsonData.AccessToken)
			log.Debug().Err(err).Interface("do", do).Msg("delete access token")
			c.JSON(200, gin.H{"api_token": auth.UserRepo.GenerateApiToken(user, false), "expire_seconds": auth.env.Config().AccessTokenLifetime})
			return
		}
	}

	c.JSON(400, gin.H{"error": "bad request"})
}
