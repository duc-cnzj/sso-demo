package integrations_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sso/app/models"
	"sso/tests"
	"testing"
)

func TestAuthController_token_index(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user, token := tests.NewUserWithToken(nil)
		user2, _ := tests.NewUserWithToken(&models.User{
			UserName: "aaa",
			Email:    "a@q.c",
		})
		repos.UserRepo.GenerateApiToken(user)
		apiToken2 := repos.UserRepo.GenerateApiToken(user2)
		w := tests.GetJson(fmt.Sprintf("/api/admin/users/%d/api_tokens", user2.ID), map[string]string{"page": "1", "page_size": "15"}, token)
		assert.Contains(t, w.Body.String(), apiToken2)
		assert.Contains(t, w.Body.String(), fmt.Sprintf("%d", user2.ID))
		assert.NotContains(t, w.Body.String(), fmt.Sprintf("%d", user.ID))
	})
}

func TestAuthController_token_index_2(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		user, token := tests.NewUserWithToken(nil)
		user2, _ := tests.NewUserWithToken(&models.User{
			UserName: "aaa",
			Email:    "a@q.c",
		})
		repos.UserRepo.GenerateApiToken(user)
		apiToken2 := repos.UserRepo.GenerateApiToken(user2)
		w := tests.GetJson(fmt.Sprintf("/api/admin/api_tokens"), map[string]string{"page": "1", "page_size": "15", "user_id": fmt.Sprintf("%d", user2.ID)}, token)
		t.Log(w.Body.String())
		assert.Contains(t, w.Body.String(), apiToken2)
		assert.Contains(t, w.Body.String(), fmt.Sprintf("%d", user2.ID))
		assert.NotContains(t, w.Body.String(), fmt.Sprintf("%d", user.ID))
	})
}
