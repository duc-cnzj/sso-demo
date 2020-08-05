package auth

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	//"net/http"
	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/jwt"
	"strings"
)

func AuthMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if strings.HasPrefix(token, "Bearer") || strings.HasPrefix(token, "bearer") {
			jwtToken := strings.TrimSpace(token[6:])
			if jwtToken != "" {
				if claims, b := jwt.ParseJwt(jwtToken); b {
					c.Set("user", models.User{}.FindByEmail(claims.Id, env))
					return
				}
			}

			c.AbortWithStatusJSON(401, gin.H{"code": 401, "error": "unauthorized!"})
		}
	}
}

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
					c.AbortWithStatusJSON(422, gin.H{"code": 422})
					return
				}
				c.Redirect(302, redirectUrl+"?access_token="+token)
				return
			}
		}

		c.Next()
	}
}
