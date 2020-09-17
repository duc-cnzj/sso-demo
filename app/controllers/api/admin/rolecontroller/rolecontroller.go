package rolecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"math"
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
	"sso/utils/form"
	"strconv"
	"strings"
)

type QueryInput struct {
	Name string `form:"name"`
	Sort string `form:"sort"`

	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type RoleStoreInput struct {
	Text          string `form:"text" json:"text" binding:"required"`
	Name          string `form:"name" json:"name" binding:"required,alpha"`
	PermissionIds []uint `form:"permission_ids" json:"permission_ids"`
}

type RoleUpdateInput struct {
	Text          string `form:"text" json:"text" binding:"required"`
	Name          string `form:"name" json:"name" binding:"required,alpha"`
	PermissionIds []uint `form:"permission_ids" json:"permission_ids"`
}

type RoleController struct {
	env *env.Env
	*api.AllRepo
}

func NewRoleController(env *env.Env) *RoleController {
	return &RoleController{env: env, AllRepo: api.NewAllRepo(env)}
}

func (role *RoleController) Index(ctx *gin.Context) {
	var (
		query QueryInput
		count int
	)

	if err := ctx.ShouldBindQuery(&query); err != nil {
		exception.ValidateException(ctx, err, role.env)
		return
	}
	log.Debug().Interface("query", query).Msg("RoleController.Index")

	var roles []models.Role
	if query.PageSize <= 0 {
		query.PageSize = 15
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	switch strings.ToLower(query.Sort) {
	case "asc":
		query.Sort = "ASC"
	case "":
		fallthrough
	case "desc":
		fallthrough
	default:
		query.Sort = "DESC"
	}

	q := role.env.GetDB().Model(&roles).Preload("Permissions")
	if query.Name != "" {
		q = q.Where("name like ?", "%"+query.Name+"%")
	}

	offset := int(math.Max(float64((query.Page-1)*query.PageSize), 0))
	q.
		Offset(offset).
		Limit(query.PageSize).
		Order("id " + query.Sort).
		Find(&roles)

	if len(roles) < query.PageSize {
		count = query.PageSize*(query.Page-1) + len(roles)
	} else {
		countQuery := role.env.GetDB().Model(&models.Role{})
		if query.Name != "" {
			countQuery = countQuery.Where("name like ?", "%"+query.Name+"%")
		}
		countQuery.Count(&count)
	}

	ctx.JSON(200, gin.H{"code": 200, "data": roles, "page": query.Page, "page_size": query.PageSize, "total": count})
}

func (role *RoleController) Store(ctx *gin.Context) {
	input := &RoleStoreInput{}
	if err := ctx.ShouldBind(input); err != nil {
		exception.ValidateException(ctx, err, role.env)
		log.Debug().Err(err).Msg("RoleController.Store")
		return
	}

	name, _ := role.RoleRepo.FindByName(input.Name)
	if name != nil {
		var errors = form.ValidateErrors{
			form.ValidateError{
				Field: "name",
				Msg:   "role exists",
			},
		}
		exception.ValidateException(ctx, errors, role.env)
		return
	}

	r := &models.Role{
		Name: input.Name,
		Text: input.Text,
	}

	err := role.RoleRepo.CreateWithPermissionIds(r, input.PermissionIds)
	if err != nil {
		log.Error().Err(err).Msg("role.RoleRepo.CreateWithPermissionIds")
		ctx.AbortWithError(500, err)
		return
	}

	permissions, _ := role.RoleRepo.FindByIdWithPermissions(r.ID)

	ctx.JSON(201, gin.H{"code": 201, "data": permissions})
}

func (role *RoleController) Show(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("role"))
	if err != nil {
		log.Panic().Err(err).Msg("RoleController.Show")

		return
	}

	r, _ := role.RoleRepo.FindByIdWithPermissions(uint(id))
	if r == nil {
		exception.ModelNotFound(ctx, "role")
		return
	}

	ctx.JSON(200, gin.H{"data": r})
}

func (role *RoleController) Update(ctx *gin.Context) {
	var input RoleUpdateInput
	id, err := strconv.Atoi(ctx.Param("role"))
	if err != nil {
		log.Panic().Err(err).Msg("RoleController.Update")

		return
	}
	if err := ctx.ShouldBind(&input); err != nil {
		exception.ValidateException(ctx, err, role.env)
		log.Debug().Err(err).Msg("RoleController.Update")

		return
	}

	r, _ := role.RoleRepo.FindById(uint(id))
	if r == nil {
		exception.ModelNotFound(ctx, "role")
		return
	}

	if input.Name != "" {
		hasRole, _ := role.RoleRepo.FindByName(input.Name)

		if hasRole != nil && hasRole.ID != r.ID {
			var errors = form.ValidateErrors{
				form.ValidateError{
					Field: "name",
					Msg:   "name exists",
				},
			}
			exception.ValidateException(ctx, errors, role.env)
			return
		}
		log.Debug().Interface("input", input).Msg("RoleController.Update")
	}

	role.env.DBTransaction(func(tx *gorm.DB) error {
		if e := tx.Model(r).Updates(map[string]interface{}{"name": input.Name, "text": input.Text}).Error; e != nil {
			log.Error().Err(e).Msg("RoleController.Update")
			return e
		}

		log.Debug().Interface("input.PermissionIds", input.PermissionIds).Msg("RoleController.Update")

		if err := role.RoleRepo.SyncPermissions(r, input.PermissionIds, tx); err != nil {
			log.Error().Err(err).Msg("RoleController.Update")
			return err
		}

		return nil
	})

	permissions, _ := role.RoleRepo.FindByIdWithPermissions(r.ID)
	ctx.JSON(200, gin.H{"code": 200, "data": permissions})
}

func (role *RoleController) Destroy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("role"))
	if err != nil {
		log.Panic().Err(err).Msg("RoleController Destroy err: ")
		return
	}

	r, _ := role.RoleRepo.FindById(uint(id))
	if r == nil {
		exception.ModelNotFound(ctx, "role")
		return
	}

	role.env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Delete(r).Error != nil {
			return tx.Delete(r).Error
		}
		if tx.Model(r).Association("Permissions").Clear().Error != nil {
			return tx.Model(r).Association("Permissions").Clear().Error
		}

		return nil
	})

	ctx.JSON(204, nil)
}

func (role *RoleController) All(c *gin.Context) {
	type res struct {
		ID   uint   `json:"id"`
		Text string `json:"text"`
	}
	var roles []*models.Role
	var result []res

	role.env.GetDB().Select([]string{"id", "text"}).Find(&roles)
	for _, v := range roles {
		result = append(result, res{
			ID:   v.ID,
			Text: v.Text,
		})
	}

	c.JSON(200, gin.H{"data": result})
}
