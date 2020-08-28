package permissioncontroller_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/server"
	"sso/tests"
	"strconv"
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

func TestPing(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		s.Engine().ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, `{"success":true}`, w.Body.String())
	})
}

func TestPermissionController_Index(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		p := &models.Permission{
			Name:    "create",
			Project: "sso-with-test",
		}
		repos.PermRepo.Create(p)

		repos.PermRepo.Create(&models.Permission{
			Name:    "show",
			Project: "sso-with-test",
		})

		data := []struct {
			name string
			token string
			data map[string]string
			code int
			res string
		}{
			{
				name:  "401",
				token: "",
				code:  401,
			},
			{
				name:  "query name",
				token: token,
				data: map[string]string{
					"name": "create-update",
				},
				code:  200,
				res: `"total":0`,
			},
			{
				name:  "query project",
				token: token,
				data: map[string]string{
					"project": "sso-with-test",
				},
				code:  200,
				res: `"total":2`,
			},
			{
				name:  "no query",
				token: token,
				code:  200,
				res: "show",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.GetJson("/api/admin/permissions",  test.data, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_Store(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		data := []struct {
			name string
			token string
			data map[string]string
			code int
			res string
		}{
			{
				name:  "401",
				token: "",
				data:  map[string]string{
					"name":   "view",
					"projec": "sso",
				},
				code:  401,
			},
			{
				name:  "success",
				token: token,
				data:  map[string]string{
					"name":   "view",
					"projec": "sso",
				},
				code:  201,
			},
			{
				name:  "permission exists",
				token: token,
				data:  map[string]string{
					"name":   "view",
					"projec": "sso",
				},
				code:  422,
				res: "permission exists",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PostJson("/api/admin/permissions", test.data, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_Show(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		p := &models.Permission{
			Name:    "create",
			Project: "sso",
		}
		repos.PermRepo.Create(p)
		id := strconv.Itoa(int(p.ID))
		data := []struct {
			name string
			token string
			data map[string]string
			code int
			res string
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
				w = tests.GetJson("/api/admin/permissions/" + id , nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_Update(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		p := &models.Permission{
			Name:    "create",
			Project: "sso",
		}
		repos.PermRepo.Create(p)

		repos.PermRepo.Create(&models.Permission{
			Name:    "show",
			Project: "sso",
		})
		id := strconv.Itoa(int(p.ID))
		data := []struct {
			name string
			token string
			data map[string]string
			code int
			res string
		}{
			{
				name:  "401",
				token: "",
				data: map[string]string{
					"name": "create-update",
					"project": "sso-big",
				},
				code:  401,
			},
			{
				name:  "success",
				token: token,
				data: map[string]string{
					"name": "create-update",
					"project": "sso-big",
				},
				code:  200,
			},
			{
				name:  "success",
				token: token,
				data: map[string]string{
					"name": "show",
					"project": "sso-big",
				},
				code:  422,
				res: "name exists",
			},
		}

		var w *httptest.ResponseRecorder
		for _, test := range data {
			t.Run(test.name, func(t *testing.T) {
				w = tests.PutJson("/api/admin/permissions/" + id ,  test.data, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_Destroy(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		p := &models.Permission{
			Name:    "create",
			Project: "sso",
		}
		repos.PermRepo.Create(p)

		id := strconv.Itoa(int(p.ID))
		data := []struct {
			name string
			token string
			code int
			res string
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
				w = tests.DeleteJson("/api/admin/permissions/" + id, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_GetByGroups(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		repos.PermRepo.Create(&models.Permission{
			Name:    "create",
			Project: "sso",
		})

		repos.PermRepo.Create(&models.Permission{
			Name:    "show",
			Project: "sso",
		})

		data := []struct {
			name string
			token string
			code int
			res string
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
				w = tests.GetJson("/api/admin/permissions_by_group", nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}

func TestPermissionController_GetPermissionProjects(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		_, token := tests.NewUserWithToken(nil)

		repos.PermRepo.Create(&models.Permission{
			Name:    "create",
			Project: "sso",
		})

		repos.PermRepo.Create(&models.Permission{
			Name:    "show",
			Project: "sso",
		})

		data := []struct {
			name string
			token string
			code int
			res string
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
				w = tests.GetJson("/api/admin/get_permission_projects", nil, test.token)
				assert.Equal(t, test.code, w.Code)
				assert.Contains(t, w.Body.String(), test.res)
				t.Log(w.Body.String())
			})
		}
	})
}