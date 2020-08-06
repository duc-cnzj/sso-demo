package auth

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sso/app/models"
	"sso/config/env"
)

func SessionMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			var u *models.User
			err := json.Unmarshal([]byte(user.(string)), &u)
			if err == nil {
				c.Set("user", u)
				c.Next()
				return
			}
		}

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
		user := session.Get("user")
		if user != nil {
			var u *models.User
			err := json.Unmarshal([]byte(user.(string)), &u)
			if err == nil {
				token := u.GenerateAccessToken(env)
				redirectUrl := c.Query("redirect_url")
				if redirectUrl == "" {
					c.Redirect(302, "/auth/select_system")
					return
				}
				c.Redirect(302, redirectUrl+"?access_token="+token)
				return
			}
		}

		c.Next()
	}
}
