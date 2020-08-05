package authcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sso/app/http/middlewares/i18n"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/form"
)

type LoginForm struct {
	UserName    string `form:"email" binding:"required"`
	Password    string `form:"password" binding:"required"`
	RedirectUrl string `form:"redirect_url" binding:"required"`
}

func New(env *env.Env) *authController {
	return &authController{env: env}
}

type authController struct {
	env *env.Env
}

func (*authController) LoginForm(ctx *gin.Context) {
	redirectUrl := ctx.Query("redirect_url")
	ctx.HTML(http.StatusOK, "login.tmpl", struct {
		RedirectUrl string
	}{
		RedirectUrl: redirectUrl,
	})
}

func (auth *authController) Login(ctx *gin.Context) {
	var loginForm LoginForm

	if err := ctx.ShouldBind(&loginForm); err != nil {
		errors := err.(validator.ValidationErrors)
		value, _ := ctx.Get(i18n.UserPreferLangKey)
		trans, _ := auth.env.GetUniversalTranslator().GetTranslator(value.(string))
		ctx.AbortWithStatusJSON(422, gin.H{"code": 422, "error": form.ErrorsToMap(errors, trans)})

		return
	}

	user := models.User{}.FindByEmail(loginForm.UserName, auth.env)
	if user == nil {
		ctx.JSON(404, gin.H{"code": 404, "error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		ctx.JSON(401, gin.H{"code": 401, "error": "email or password error!"})

		return
	}

	token := user.GenerateAccessToken(auth.env)

	ctx.Redirect(200,  ,loginForm.RedirectUrl+"?access_token="+token)
}

func (auth *authController) Me(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(*models.User)

	c.JSON(200, gin.H{"data": user})
}
