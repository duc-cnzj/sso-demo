package filters

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"sso/pkg/filters"
)

type ApiTokenInput struct {
	UserId uint `json:"user_id" form:"user_id" uri:"user"`
}

func NewApiTokenFilter(ctx *gin.Context) (filters.Filterable, error) {
	var input ApiTokenInput

	if err := ctx.ShouldBindUri(&input); err != nil {
		return nil, err
	}

	if err := ctx.ShouldBind(&input); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"user_id": input.UserId,
	}
	log.Debug().Interface("ad", data).Msg("filter")
	f := filters.NewFilter(data)

	f.RegisterFunc("user_id", UserId)

	return f, nil
}

func UserId(f filters.Filterable) filters.GormScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		var (
			userId interface{}
			err    error
		)
		if userId, err = f.GetInput("user_id"); err != nil {
			return db
		}

		id := userId.(uint)
		if id > 0 {
			return db.Where("user_id = ?", id)
		}

		return db
	}
}
