package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
	"sso/repositories/user_repository"
)

const HttpAuthToken = "X-Request-Token"

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

		c.Redirect(302, "/login")
		c.Abort()
	}
}

func GuestMiddleware(env *env.Env) gin.HandlerFunc {
	userRepo := user_repository.NewUserRepository(env)

	return func(c *gin.Context) {
		session := sessions.Default(c)
		u, ok := session.Get("user").(*models.User)
		if ok {
			if !CheckLogoutTokenIsChanged(u.LogoutToken, u.ID, env) {
				token := userRepo.GenerateAccessToken(u)
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
	userRepo := user_repository.NewUserRepository(env)

	user, _ := userRepo.FindById(id)
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
	userRepo := user_repository.NewUserRepository(env)

	// todo 需要优化，这个接口访问最频繁，但是查了2次数据库
	return func(c *gin.Context) {
		token := c.Request.Header.Get(HttpAuthToken)
		if token != "" {
			user, err := userRepo.FindByToken(token, true)
			log.Debug().Interface("user", user).Err(err).Msg("ApiMiddleware")
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
