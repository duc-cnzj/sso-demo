package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
	"sso/repositories/user_repository"
	"strconv"
	"strings"
	"time"
)

type SsoJwtClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

func AuthMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ssoClaims SsoJwtClaims
		token := GetBearerToken(c)

		if token != "" && !KeyInBlacklist(token, env) {
			t, err := jwt.ParseWithClaims(token, &ssoClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte(env.Config().JwtSecret), nil
			})

			if err != nil {
				log.Error().Err(err).Msg("jwt.AuthMiddleware: " + token)
				Unauthorized(c)
				return
			}

			if claims, ok := t.Claims.(*SsoJwtClaims); ok && t.Valid {
				log.Debug().Interface("t", claims).Msg("jwt.AuthMiddleware")
				userRepo := user_repository.NewUserRepository(env)
				pretty, _ := userRepo.LoadUserRoleAndPermissionPretty(claims.User, "sso")
				claims.User.CurrentRoles = pretty.Roles
				claims.User.CurrentRoles = pretty.Permissions

				env.Auth().SetUser(claims.User)
				c.Set("user", claims.User)
				c.Next()
				return
			}
		}

		Unauthorized(c)
	}
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(401, gin.H{
		"code": 401,
		"msg":  "Unauthorized!",
	})
}

func GetBearerToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")

	if token != "" && (strings.HasPrefix(token, "Bearer") || strings.HasPrefix(token, "bearer")) {
		return strings.TrimSpace(token[6:])
	}

	log.Debug().Msg("GetBearerToken:" + token)
	return ""
}

func GenerateToken(user *models.User, env *env.Env) (string, error) {
	mySigningKey := []byte(env.Config().JwtSecret)

	uId := strconv.Itoa(int(user.ID))
	exp := time.Second * time.Duration(env.Config().JwtExpiresSeconds)

	claims := &SsoJwtClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
			Id:        uId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(mySigningKey)
}

func Parse(token string, env *env.Env) (string, error) {
	var ssoClaims SsoJwtClaims
	if token != "" {
		fmt.Println(token)
		t, err := jwt.ParseWithClaims(token, &ssoClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.Config().JwtSecret), nil
		})

		if err != nil {
			return "", err
		}

		if claims, ok := t.Claims.(*SsoJwtClaims); ok && t.Valid {
			log.Debug().Interface("t", claims).Msg("jwt.AuthMiddleware")
			marshal, _ := json.MarshalIndent(claims, "", "\t")
			return string(marshal), nil
		}
	}

	return "", errors.New("MissToken")
}
