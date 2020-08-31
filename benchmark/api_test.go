package api_test

import (
	"os"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/server"
	"sso/tests"
	"testing"
)

var (
	repos *api.AllRepo
	s     *server.Server
	token string
)

func TestMain(m *testing.M) {
	pwd, _ := os.Getwd()

	s, repos = tests.MainHelper(pwd + "/../.env.testing")
	begin := s.Env().GetDB().Begin()
	s.Env().SetDB(begin)
	_, token = tests.NewUserWithToken(&models.User{
		UserName: "a111dsad",
		Email:    "adad@a.c",
		Password: "123",
	})
	code := m.Run()
	s.Env().GetDB().Rollback()
	os.Exit(code)
}

func BenchmarkPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tests.GetJson("/ping", nil, "")
	}
}

func BenchmarkNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tests.GetJson("/not_found", nil, "")
	}
}
