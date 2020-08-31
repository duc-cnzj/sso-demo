package integrations_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sso/app/models"
	"sso/tests"
	"strconv"
	"testing"
)

func TestRoleController_Store(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		r := &models.Role{
			Name: "role one",
		}
		repos.RoleRepo.Create(r)
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
					"name": "role",
				},
				code: 401,
			},
			{
				name: "success",
				body: map[string]string{
					"name": "role",
				},
				token: token,
				code:  201,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PostJson("/api/admin/roles", test.body, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestRoleController_Show(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		r := &models.Role{
			Name: "",
		}
		repos.RoleRepo.Create(r)
		id := strconv.Itoa(int(r.ID))
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
				w = tests.GetJson("/api/admin/roles/"+id, nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestRoleController_Update(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		r := &models.Role{
			Name: "sso",
		}
		repos.RoleRepo.Create(r)
		id := strconv.Itoa(int(r.ID))
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
					"name": "update",
				},
				code: 401,
			},
			{
				name:  "success",
				token: token,
				body: map[string]string{
					"name": "update",
				},
				code: 200,
				res:  "update",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PutJson("/api/admin/roles/"+id, test.body, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestRoleController_Destroy(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		r := &models.Role{
			Name: "",
		}
		repos.RoleRepo.Create(r)
		id := strconv.Itoa(int(r.ID))
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
				w = tests.DeleteJson("/api/admin/roles/"+id, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestRoleController_All(t *testing.T) {
	type Result struct {
		Data []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}
	}
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		i := 10
		for i > 0 {
			r := &models.Role{
				Name: "sso_" + strconv.Itoa(i),
			}
			repos.RoleRepo.Create(r)
			i--
		}

		data := []struct {
			name  string
			token string
			code  int
			count int
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
				count: 10,
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.GetJson("/api/admin/all_roles", nil, test.token)
				assert.Equal(t, test.code, w.Code)
				var resultData Result
				json.Unmarshal(w.Body.Bytes(), &resultData)
				assert.Equal(t, test.count, len(resultData.Data))
			})
		}
	})
}
