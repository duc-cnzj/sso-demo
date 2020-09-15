package integrations_test

import (
	"github.com/stretchr/testify/assert"
	"sso/app/models"
	"sso/tests"
	"testing"
)

func TestAuthController_Info(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user)
		w := tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
		t.Log(w.Body.String())
	})
}

func TestAuthController_login(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		var apiToken models.ApiToken
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user)
		s.Env().GetDB().Where("api_token = ?", token).First(&apiToken)
		assert.Nil(t, apiToken.LastUseAt)
		w := tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
		s.Env().GetDB().Where("api_token = ?", token).First(&apiToken)
		assert.NotNil(t, apiToken.LastUseAt)
	})
}

// 用户可以有多个api_token，只要不过期就可以获取信息
func TestAuthController_login_use_diff_effective_token(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user)
		w := tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
		repos.UserRepo.GenerateApiToken(user)
		w = tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
	})
}

// 生成2个api_token，当用户登出时，统统失效
func TestAuthController_Logout(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user := tests.NewUser(nil)
		token := repos.UserRepo.GenerateApiToken(user)
		w := tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 200, w.Code)
		token2 := repos.UserRepo.GenerateApiToken(user)

		w = tests.WebPostJson("/api/user/info", nil, token2)
		assert.Equal(t, 200, w.Code)

		w = tests.WebPostJson("/api/logout", nil, token)
		assert.Equal(t, 204, w.Code)
		id, _ := repos.UserRepo.FindById(user.ID)
		assert.NotEqual(t, user.LogoutToken, id.LogoutToken)
		w = tests.WebPostJson("/api/user/info", nil, token2)
		assert.Equal(t, 401, w.Code)
		w = tests.WebPostJson("/api/user/info", nil, token)
		assert.Equal(t, 401, w.Code)
	})
}
