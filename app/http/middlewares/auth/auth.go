package auth

import (
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
