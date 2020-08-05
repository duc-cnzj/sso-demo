package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type User interface {
	GetJwtID() string
}

var expiredTime = time.Hour * 24
var mySigningKey = "secret"

// 生成令牌  创建jwt风格的token
func GenerateToken(user User) string {
	claims := &jwt.StandardClaims{
		Audience:  "uco users",
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Id:        user.GetJwtID(),
		Issuer:    "uco",
		Subject:   "sso",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		log.Fatal("jwt signed string error", err)
	}

	return tokenString
}

func ParseJwt(tokenStr string) (*jwt.StandardClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		log.Println(err)
		return nil, false
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, true
	}
	log.Println(token)

	return nil, false
}
