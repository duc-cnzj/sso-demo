package authcontroller_test

import (
	"github.com/stretchr/testify/assert"
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

	s, repos = tests.MainHelper(pwd + "/../../../../../.env.testing")

	os.Exit(m.Run())
}

func TestAuthController_Info(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user, true)
		w := tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
		t.Log(w.Body.String())
	})
}

func TestAuthController_Logout(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user, true)
		w := tests.WebPostJson("/api/logout", nil, token)
		assert.Equal(t, 204, w.Code)
		id, _ := repos.UserRepo.FindById(user.ID)
		assert.NotEqual(t, user.LogoutToken, id.LogoutToken)
		assert.NotEqual(t, user.ApiToken.String, id.ApiToken.String)
	})
}
