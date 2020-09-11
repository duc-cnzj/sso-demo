package integrations_test

import (
	"github.com/coreos/etcd/pkg/testutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso/app/filters"
	"testing"
)

func NewFilterCtx(data map[string]string) *gin.Context {
	req, _ := http.NewRequest("GET", "", nil)
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")

	return &gin.Context{
		Request: req,
	}
}

//UserName string `form:"user_name" json:"user_name"`
//Email    string `form:"email" json:"email"`
//
//Page     int    `form:"page" json:"page"`
//PageSize int    `form:"page_size" json:"page_size"`
//Sort     string `form:"sort" json:"sort"`
func TestUserFilter(t *testing.T) {
	ctx := NewFilterCtx(map[string]string{
		"user_name": "duc",
		"email":     "duc@qq.com",
		"sort":      "asc",
	})
	filter, _ := filters.NewUserFilter(ctx)
	//var res models.User
	//s.Env().GetDB().LogMode(true)
	//s.Env().GetDB().Scopes(filter.Apply()...).Find(&res)
	//s.Env().GetDB().Scopes(filter.Apply()...).Find(&res)
	testutil.AssertEqual(t, len(filter.Apply()), len(filter.Apply()))
}

func BenchmarkFilter(b *testing.B) {
	ctx := NewFilterCtx(map[string]string{
		"user_name": "duc",
		"email":     "duc@qq.com",
		"sort":      "asc",
	})
	filter, _ := filters.NewUserFilter(ctx)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply()
	}
}
