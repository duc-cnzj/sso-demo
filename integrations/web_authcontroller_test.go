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
		role := &models.Role{
			Name: "admin",
		}
		p1 := &models.Permission{
			Name:    "view",
			Project: "sso",
		}
		p2 := &models.Permission{
			Name:    "new",
			Project: "sso",
		}
		p3 := &models.Permission{
			Name:    "view",
			Project: "micro",
		}
		p4 := &models.Permission{
			Name:    "delete",
			Project: "micro",
		}
		p5 := &models.Permission{
			Name:    "delete",
			Project: "micro",
		}
		repos.RoleRepo.Create(role)
		repos.PermRepo.Create(p1)
		repos.PermRepo.Create(p2)
		repos.PermRepo.Create(p3)
		repos.PermRepo.Create(p4)
		repos.PermRepo.Create(p5)
		assert.Equal(t, p5.ID, p4.ID)
		repos.RoleRepo.SyncPermissions(role, []uint{p1.ID, p2.ID, p3.ID, p4.ID}, nil)
		token := repos.UserRepo.GenerateApiToken(user)
		repos.UserRepo.SyncRoles(user, []*models.Role{role})
		w := tests.WebPostJson("/api/user/info/projects/micro", nil, token)
		assert.Equal(t, 200, w.Code)
		t.Log(w.Body.String())
		w = tests.WebPostJson("/api/user/info", nil, token)
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
