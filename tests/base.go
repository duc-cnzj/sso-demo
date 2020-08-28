package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"log"
	"sso/app/http/controllers/api"
	"sso/app/http/middlewares/jwt"
	"sso/app/models"
	"sso/server"
)

var (
	repos *api.AllRepo
	s     *server.Server
)

func NewTestServer(path string) (*server.Server, error) {
	var s = &server.Server{}
	if err := s.Init(path, ""); err != nil {
		return nil, err
	}

	return s, nil
}

func MainHelper(env string) (*server.Server, *api.AllRepo) {
	var (
		err error
	)

	zerolog.SetGlobalLevel(zerolog.Disabled)
	gin.SetMode(gin.ReleaseMode)
	s, err = NewTestServer(env)
	if err != nil {
		log.Panic(err)
	}

	repos = api.NewAllRepo(s.Env())
	return s, repos
}

func WarpTxRollback(s *server.Server, fn func()) {
	db := s.Env().GetDB()
	s.Env().DBTransaction(func(tx *gorm.DB) error {
		s.Env().SetDB(tx)
		fn()
		tx.Rollback()
		return nil
	})
	s.Env().SetDB(db)
}

func NewUserWithToken(user *models.User) (*models.User, string) {
	u := NewUser(user)

	generateToken, _ := jwt.GenerateToken(u, s.Env())

	return u, generateToken
}

func NewUser(user *models.User) *models.User {
	pwd, _ := repos.UserRepo.GeneratePwd("12345")
	var u *models.User
	if user != nil {
		u = user
	} else {
		u = &models.User{
			UserName: "duc",
			Email:    "duc@duc.com",
			Password: pwd,
		}
	}
	if err := repos.UserRepo.Create(u); err != nil {
		return nil
	}

	return u
}
