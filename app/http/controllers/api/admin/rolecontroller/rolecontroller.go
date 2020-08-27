package rolecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"math"
	"sso/app/http/controllers/api"
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
	Name          string `form:"name" json:"name"`
	PermissionIds []uint `form:"permission_ids" json:"permission_ids"`
}

type RoleUpdateInput struct {
	Name          string `form:"name" json:"name"`
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
	}

	e := role.env.DBTransaction(func(tx *gorm.DB) error {
		if err := tx.Create(r).Error; err != nil {
			return err
		}

		if input.PermissionIds != nil {
			permissionByIds, _ := role.PermRepo.FindByIds(input.PermissionIds)
			if err := tx.Model(r).Association("Permissions").Clear().Error; err != nil {
				return err
			}

			if err := tx.Model(r).Association("Permissions").Append(toInterfaceSlice(permissionByIds)...).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if e != nil {
		log.Panic().Err(e).Msg("RoleController.Store")
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
		e := role.env.GetDB().Model(r).Update("name", input.Name)
		log.Debug().Err(e.Error).Msg("RoleController.Update")
	}

	log.Debug().Interface("input.PermissionIds", input.PermissionIds).Msg("RoleController.Update")

	if input.PermissionIds != nil {
		ps, _ := role.PermRepo.FindByIds(input.PermissionIds)
		role.env.GetDB().Model(r).Association("Permissions").Clear()
		role.env.GetDB().Model(r).Association("Permissions").Append(toInterfaceSlice(ps)...)
	}

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
		Name string `json:"name"`
	}
	var roles []*models.Role
	var result []res

	role.env.GetDB().Select([]string{"name", "id"}).Find(&roles)
	for _, v := range roles {
		result = append(result, res{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	c.JSON(200, gin.H{"data": result})
}

func toInterfaceSlice(slice interface{}) []interface{} {
	permissions := slice.([]*models.Permission)
	newS := make([]interface{}, len(permissions))
	for i, v := range permissions {
		newS[i] = v
	}

	return newS
}
