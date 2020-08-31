package integrations_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sso/app/middlewares/jwt"
	"sso/app/models"
	"sso/tests"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		type LoginForm struct {
			UserName string `json:"email"`
			Password string `json:"password"`
		}
		body := &bytes.Buffer{}
		json.NewEncoder(body).Encode(LoginForm{
			UserName: "1@q.c",
			Password: "1234",
		})
		v := url.Values{}
		v.Add("email", "1@q.c")
		v.Add("password", "1234")
		jack := url.Values{}
		jack.Add("email", "jack@qq.com")
		jack.Add("password", "12345")

		data := []struct {
			name        string
			body        io.Reader
			contentType string
			want        struct {
				code int
				res  string
			}
		}{
			{
				name:        "login use json input(fail login)",
				body:        body,
				contentType: "application/json",
				want: struct {
					code int
					res  string
				}{code: 401, res: `{"code":401,"msg":"Unauthorized!"}`},
			},
			{
				name:        "login use form input(fail login)",
				body:        strings.NewReader(v.Encode()),
				contentType: "application/x-www-form-urlencoded",
				want: struct {
					code int
					res  string
				}{code: 401, res: `{"code":401,"msg":"Unauthorized!"}`},
			},
			{
				name:        "login use form input(success login)",
				body:        strings.NewReader(jack.Encode()),
				contentType: "application/x-www-form-urlencoded",
				want: struct {
					code int
					res  string
				}{
					code: 200,
					res:  "lifetime",
				},
			},
		}

		pwd, _ := repos.UserRepo.GeneratePwd("12345")

		s.Env().GetDB().Create(&models.User{
			UserName: "jack",
			Email:    "jack@qq.com",
			Password: pwd,
		})

		for _, tt := range data {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/admin/login", tt.body)
				req.Header.Add("Content-Type", tt.contentType)
				s.Engine().ServeHTTP(w, req)
				assert.Equal(t, tt.want.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.want.res)
			})
		}
	})
}

func TestLogout(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/admin/logout", nil)
		s.Engine().ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		user := tests.NewUser(nil)
		token, _ := jwt.GenerateToken(user, s.Env())
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/admin/logout", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		s.Engine().ServeHTTP(w, req)
		assert.Equal(t, 204, w.Code)
		assert.True(t, jwt.KeyInBlacklist(token, s.Env()))
		get := s.Env().RedisPool().Get()
		do, _ := get.Do("FLUSHALL")
		t.Log(do)
		get.Close()
		assert.False(t, jwt.KeyInBlacklist(token, s.Env()))
	})
}

func TestInfo(t *testing.T) {
	tests.WarpTxRollback(s, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/admin/user/info", nil)
		s.Engine().ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		user := tests.NewUser(nil)
		token, _ := jwt.GenerateToken(user, s.Env())

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/admin/user/info", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		s.Engine().ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), user.Email)
	})
}
