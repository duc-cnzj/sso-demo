package authcontroller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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
	return &authController{env: env}
}

type authController struct {
	env *env.Env
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

	if err := ctx.ShouldBind(&loginForm); err != nil {
		exception.ValidateException(ctx, err, auth.env)

		return
	}

	user := models.User{}.FindByEmail(loginForm.UserName, auth.env)
	printErrorBack := func() {
		ctx.HTML(200, "login.tmpl", LoginFormVal{
			Errors: []string{"username or password error."},
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
	err := session.Save()
	log.Println(err)

	user.UpdateLastLoginAt(auth.env)

	if loginForm.RedirectUrl == "" {
		ctx.Redirect(302, "/auth/select_system")
		return
	}

	token := user.GenerateAccessToken(auth.env)

	ctx.Redirect(302, loginForm.RedirectUrl+"?access_token="+token)
}

func (auth *authController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user, ok := session.Get("user").(*models.User)
	session.Clear()
	session.Save()
	if ok {
		// 让用户的api调用不能再使用
		user.GenerateApiToken(auth.env, true)
		// 登出sso系统
		user.GenerateLogoutToken(auth.env)
	}

	c.Redirect(302, "/login")
}

func (auth *authController) SelectSystem(c *gin.Context) {
	c.HTML(200, "select_system.tmpl", nil)
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

	log.Println(jsonData.AccessToken)
	conn := auth.env.RedisPool().Get()

	defer conn.Close()

	id, err := redis.Int(conn.Do("GET", jsonData.AccessToken))
	log.Println(id, err)
	if err == nil {
		user := models.User{}.FindById(uint(id), auth.env)
		if user != nil {
			do, err := conn.Do("DEL", jsonData.AccessToken)
			log.Println("delete access token", do, err)
			c.JSON(200, gin.H{"api_token": user.GenerateApiToken(auth.env, false), "expire_seconds": auth.env.Config().AccessTokenLifetime})
			return
		}
	}

	c.JSON(400, gin.H{"error": "bad request"})
}