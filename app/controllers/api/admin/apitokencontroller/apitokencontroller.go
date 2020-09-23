package apitokencontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	"sso/app/controllers/api"
	"sso/app/filters"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
)

type apiTokenController struct {
	env *env.Env
	*api.AllRepo
}

type Uri struct {
	UserId uint
}

type Paginate struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func New(env *env.Env) *apiTokenController {
	return &apiTokenController{env: env, AllRepo: api.NewAllRepo(env)}
}

func (token *apiTokenController) Index(c *gin.Context) {
	var (
		tokens   []models.ApiToken
		paginate Paginate
		count    int
	)

	if err := c.ShouldBind(&paginate); err != nil {
		exception.ValidateException(c, err, token.env)
		return
	}

	filter, _ := filters.NewApiTokenFilter(c)
	offset := int(math.Max(float64((paginate.Page-1)*paginate.PageSize), 0))

	token.env.GetDB().
		Scopes(filter.Apply()...).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"users.id", "users.user_name"})
		}).
		Offset(offset).
		Limit(paginate.PageSize).
		Order("id DESC").
		Find(&tokens)

	if len(tokens) < paginate.PageSize {
		count = paginate.PageSize*(paginate.Page-1) + len(tokens)
	} else {
		token.env.GetDB().Model(&models.ApiToken{}).Scopes(filter.Apply()...).Count(&count)
	}
	c.JSON(200, gin.H{"code": 200, "data": tokens, "page": paginate.Page, "page_size": paginate.PageSize, "total": count})
}
