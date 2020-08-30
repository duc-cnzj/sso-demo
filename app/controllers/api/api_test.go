package api_test

import (
	"github.com/magiconair/properties/assert"
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

	s, repos = tests.MainHelper(pwd + "/../../../.env.testing")

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	w := tests.GetJson("/ping", nil, "")
	assert.Equal(t, `{"success":true}`, w.Body.String())
}

func TestNotFound(t *testing.T) {
	w := tests.GetJson("/not_found", nil, "")
	assert.Equal(t, `{"code":404,"message":"Page not found"}`, w.Body.String())
}