package integrations_test

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sso/app/models"
	"sso/tests"
	"strconv"
	"testing"
)

func TestUserController_Index(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(&models.User{
			UserName: "duc",
			Email:    "1@q.com",
		})
		data := []struct {
			name  string
			token string
			data  map[string]string
			code  int
			res   string
		}{
			{
				name:  "success",
				token: token,
				data: map[string]string{
					"user_name": "duc",
				},
				code: 200,
				res:  "duc",
			},
			{
				name:  "success",
				token: token,
				data: map[string]string{
					"user_name": "aaa",
				},
				code: 200,
				res:  "",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.GetJson("/api/admin/users", test.data, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_Store(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)
		data := []struct {
			name  string
			token string
			body  map[string]string
			code  int
			res   string
		}{
			{
				name:  "401",
				token: "",
				body: map[string]string{
					"user_name": "duc",
					"password":  "duc",
					"email":     "123@duc.com",
				},
				code: 401,
			},
			{
				name:  "success",
				token: token,
				body: map[string]string{
					"user_name": "duc",
					"password":  "duc",
					"email":     "123@duc.com",
				},
				code: 201,
			},
			{
				name:  "invalid 422",
				token: token,
				body: map[string]string{
					"user_name": "duc",
					"password":  "duc",
					"email":     "123@duc.com",
				},
				code: 422,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PostJson("/api/admin/users", test.body, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_Show(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		u, token := tests.NewUserWithToken(nil)
		atoi := strconv.Itoa(int(u.ID))
		data := []struct {
			name  string
			token string
			code  int
			res   string
		}{
			{
				name:  "401",
				token: "",
				code:  401,
			},
			{
				name:  "success",
				token: token,
				code:  200,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.GetJson("/api/admin/users/"+atoi, nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_Update(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		u, token := tests.NewUserWithToken(nil)
		atoi := strconv.Itoa(int(u.ID))
		log.Debug().Msg(atoi)
		data := []struct {
			name  string
			token string
			code  int
			body  map[string]string
			res   string
		}{
			{
				name:  "401",
				token: "",
				body: map[string]string{
					"name": "abc",
				},
				code: 401,
			},
			{
				name:  "success",
				token: token,
				body: map[string]string{
					"name":  "abc",
					"email": "abc@abc.com",
				},
				code: 200,
				res:  "abc",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PutJson("/api/admin/users/"+atoi, test.body, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_Destroy(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		u, token := tests.NewUserWithToken(nil)
		atoi := strconv.Itoa(int(u.ID))
		data := []struct {
			name  string
			token string
			code  int
			res   string
		}{
			{
				name:  "401",
				token: "",
				code:  401,
			},
			{
				name:  "success",
				token: token,
				code:  204,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.DeleteJson("/api/admin/users/"+atoi, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_SyncRoles(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		u, token := tests.NewUserWithToken(nil)
		atoi := strconv.Itoa(int(u.ID))
		r1 := &models.Role{
			Name: "1",
		}
		repos.RoleRepo.Create(r1)
		r2 := &models.Role{
			Name: "2",
		}
		repos.RoleRepo.Create(r2)
		data := []struct {
			name  string
			token string
			body  map[string]interface{}
			code  int
			res   string
		}{
			{
				name:  "401",
				token: "",
				code:  401,
			},
			{
				name:  "success",
				token: token,
				body: map[string]interface{}{
					"role_ids": []uint{r1.ID, r2.ID},
				},
				code: 200,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PostJson("/api/admin/users/"+atoi+"/sync_roles", test.body, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestUserController_ForceLogout(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		// api token logout token 都会不一样
		u, token := tests.NewUserWithToken(&models.User{
			UserName:    "duc",
			Email:       "1@q.c",
			LogoutToken: "1234",
			Password:    "123",
		})
		atoi := strconv.Itoa(int(u.ID))
		data := []struct {
			name  string
			token string
			code  int
			res   string
		}{
			{
				name:  "401",
				token: "",
				code:  401,
			},
			{
				name:  "success",
				token: token,
				code:  200,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PostJson("/api/admin/users/"+atoi+"/force_logout", nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
				if w.Code == 200 {
					byId, _ := repos.UserRepo.FindById(u.ID)
					assert.NotEqual(t, byId.LogoutToken, u.LogoutToken)
				}
			})
		}
	})

}
