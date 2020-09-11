package filters

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"sso/pkg/filters"
	"strings"
)

type UserInput struct {
	UserName string `form:"user_name" json:"user_name"`
	Email    string `form:"email" json:"email"`

	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Sort     string `form:"sort" json:"sort"`
}

type UserFilter struct {
	input     *UserInput
	scopes    []filters.GormScopeFunc
	filters   map[string]func() filters.GormScopeFunc
	ApplyFunc func(f filters.Filterable) []filters.GormScopeFunc
}

func NewUserFilter(ctx *gin.Context) (*UserFilter, error) {
	var input UserInput
	if err := ctx.ShouldBind(&input); err != nil {
		return nil, err
	}

	uf := &UserFilter{
		input:     &input,
		scopes:    make([]filters.GormScopeFunc, 0),
		ApplyFunc: filters.DefaultApply(),
	}

	uf.filters = map[string]func() filters.GormScopeFunc{
		"user_name": uf.UserName,
		"email":     uf.Email,
		"sort":      uf.Sort,
	}
	return uf, nil
}

func (f *UserFilter) All() []string {
	var all []string
	for s := range f.filters {
		all = append(all, s)
	}
	return all
}

func (f *UserFilter) Scopes() []filters.GormScopeFunc {
	return f.scopes
}

func (f *UserFilter) ResetScopes() {
	f.scopes = make([]filters.GormScopeFunc, 0)
}

func (f *UserFilter) GetFuncByName(s string) func() filters.GormScopeFunc {
	if fn, ok := f.filters[s]; ok {
		return fn
	}

	return nil
}

func (f *UserFilter) Apply() []filters.GormScopeFunc {
	f.scopes = f.ApplyFunc(f)
	return f.scopes
}

func (f *UserFilter) Push(scope filters.GormScopeFunc) {
	f.scopes = append(f.scopes, scope)
}

func (f *UserFilter) UserName() filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_name like ?", f.input.UserName)
	}
}

func (f *UserFilter) Email() filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email like ?", "%"+f.input.Email+"%")
	}
}

func (f *UserFilter) Sort() filters.GormScopeFunc {
	var sort string
	switch strings.ToLower(f.input.Sort) {
	case "asc":
		sort = "ASC"
	case "":
		fallthrough
	case "desc":
		fallthrough
	default:
		sort = "DESC"
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id " + sort)
	}
}
