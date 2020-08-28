package permissioncontroller_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"sso/app/http/controllers/api"
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

	s, repos = tests.MainHelper(pwd + "/../../../../../../.env.testing")

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

func TestStore(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		w := PostJson("/api/admin/permissions", map[string]string{
			"name":   "view",
			"projec": "sso",
		}, "")
		assert.Equal(t, 401, w.Code)

		_, token := tests.NewUserWithToken(nil)
		w = PostJson("/api/admin/permissions", map[string]string{
			"name":   "view",
			"projec": "sso",
		}, token)
		assert.Equal(t, 201, w.Code)
	})
}

func GetJson(url string, token string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}

func PostJson(url string, data interface{}, token string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	s.Engine().ServeHTTP(w, req)
	return w
}
