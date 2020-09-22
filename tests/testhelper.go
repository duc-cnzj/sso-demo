package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"net/http/httptest"
	"sso/app/controllers/api"
	webAuthMiddleware "sso/app/middlewares/auth"
	"sso/app/middlewares/jwt"
	"sso/app/models"
	"sso/server"
	"sync"
)

var (
	repos         *api.AllRepo
	s             *server.Server
	migrateModels = []interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.ApiToken{},
	}
	mu = &sync.Mutex{}
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
	migrate := s.Env().GetDB().AutoMigrate(migrateModels...)

	repos = api.NewAllRepo(s.Env())

	if migrate.Error != nil {
		log.Fatal("migrate.Error", migrate.Error.Error())
	}

	log.Println("migrate ok!")

	return s, repos
}

func WarpTxRollback(s *server.Server, fn func()) {
	mu.Lock()
	defer mu.Unlock()
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

func GetJson(url string, data map[string]string, token string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}

func PostJson(url string, data interface{}, token string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}

func PutJson(url string, data interface{}, token string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}

func DeleteJson(url string, token string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}

func WebPostJson(url string, data interface{}, token string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if token != "" {
		req.Header.Add(webAuthMiddleware.HttpAuthToken, token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}
