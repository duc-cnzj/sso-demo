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

func NewUserFilter(ctx *gin.Context) (filters.Filterable, error) {
	var input UserInput
	if err := ctx.ShouldBind(&input); err != nil {
		return nil, err
	}

	f := filters.NewFilter(map[string]interface{}{
		"email":     input.Email,
		"sort":      input.Sort,
		"user_name": input.UserName,
	})

	f.RegisterFilterFunc("user_name", UserName)
	f.RegisterFilterFunc("email", Email)
	f.RegisterFilterFunc("sort", Sort)

	return f, nil
}

func UserName(f filters.Filterable) filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		if get, err := f.Get("user_name"); err != nil {
			return db.Where("user_name like ?", get)
		}

		return db
	}
}

func Email(f filters.Filterable) filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		if get, err := f.Get("email"); err != nil {
			str := get.(string)
			return db.Where("email like ?", "%"+str+"%")
		}

		return db
	}
}

func Sort(f filters.Filterable) filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		var (
			sort string
			err  error
			get  interface{}
		)
		if get, err = f.Get("sort"); err != nil {
			return db
		}
		sort = get.(string)
		switch strings.ToLower(sort) {
		case "asc":
			sort = "ASC"
		case "":
			fallthrough
		case "desc":
			fallthrough
		default:
			sort = "DESC"
		}

		return db.Order("id " + sort)
	}
}
