package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
)

func SessionMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user, ok := session.Get("user").(*models.User)
		log.Debug().Interface("user", user).Msg("SessionMiddleware")
		if ok {
			if !CheckLogoutTokenIsChanged(user.LogoutToken, user.ID, env) {
				c.Set("user", user)
				c.Next()
				return
			} else {
				session.Clear()
				session.Save()
			}
		}

		// todo 有没有更好的办法
		Scheme := "http://"
		if c.Request.TLS != nil {
			Scheme = "https://"
		}
		c.Redirect(302, "/login?redirect_url="+Scheme+c.Request.Host+c.Request.URL.Path)
		c.Abort()
	}
}

func GuestMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		u, ok := session.Get("user").(*models.User)
		if ok {
			if !CheckLogoutTokenIsChanged(u.LogoutToken, u.ID, env) {
				token := u.GenerateAccessToken(env)
				redirectUrl := c.Query("redirect_url")
				if redirectUrl == "" {
					c.Redirect(302, "/")
					return
				}
				c.Redirect(302, redirectUrl+"?access_token="+token)
				return
			} else {
				session.Clear()
				session.Save()
			}
		}
		c.Next()
	}
}

// 如果用户做了登出操作(不管是在sso登出还是在a.com登出)，则都会改变token，导致sso登录过期
func CheckLogoutTokenIsChanged(sessionLogoutToken string, id uint, env *env.Env) bool {
	user := models.User{}.FindById(id, env)
	if user == nil {
		return true
	}
	log.Debug().Fields(map[string]interface{}{
		"sessionLogoutToken": sessionLogoutToken,
		"user.Password":      user.Password,
	}).Msg("sessionLogoutToken")

	if user.LogoutToken != "" && user.LogoutToken != sessionLogoutToken {
		return true
	}

	return false
}

func ApiMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Request-Token")
		if token != "" {
			user := models.User{}.FindByToken(token, env)
			if user != nil {
				c.Set("user", user)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "Unauthorized!",
		})
	}
}
