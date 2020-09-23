package integrations_test

import (
	"os"
	"sso/app/controllers/api"
	"sso/server"
	"sso/tests"
	"testing"
)

var (
	repos *api.AllRepo
	s     *server.Server
)

func TestMain(m *testing.M) {
	pwd, _ := os.Getwd()

	s, repos = tests.MainHelper(pwd + "/../.env.testing")
	//s.Env().GetDB().LogMode(true)

	os.Exit(m.Run())
}
