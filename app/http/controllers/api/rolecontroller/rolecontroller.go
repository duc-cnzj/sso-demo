package rolecontroller

import (
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/exception"
	"sso/utils/form"
	"strconv"
)

type QueryInput struct {
	Name string `form:"name"`

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
}

func NewRoleController(env *env.Env) *RoleController {
	return &RoleController{env: env}
}

func (role *RoleController) Index(ctx *gin.Context) {
	var query QueryInput
	if err := ctx.ShouldBindQuery(&query); err != nil {
		exception.ValidateException(ctx, err, role.env)
		return
	}
	log.Println(query)

	var roles []models.Role
	if query.PageSize <= 0 {
		query.PageSize = 15
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	q := role.env.GetDB().Model(&roles)
	if query.Name != "" {
		q = q.Where("name like ?", "%"+query.Name+"%")
	}

	offset := int(math.Max(float64((query.Page-1)*query.PageSize), 0))
	q.
		Offset(offset).
		Limit(query.PageSize).
		Order("id DESC").
		Find(&roles)

	ctx.JSON(200, gin.H{"code": 200, "data": roles})
}

func (role *RoleController) Store(ctx *gin.Context) {
	input := &RoleStoreInput{}
	if err := ctx.ShouldBind(input); err != nil {
		exception.ValidateException(ctx, err, role.env)
		log.Println("RoleController Store err: ", err)
		return
	}

	name := models.Role{}.FindByName(input.Name, role.env)
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

	newRole := role.env.GetDB().Create(r)
	if newRole.Error != nil {
		log.Panicln(newRole.Error.Error())
	}
	if input.PermissionIds != nil {
		role.env.GetDB().Model(newRole).Association("Permissions").Replace(models.Permission{}.FindByIds(input.PermissionIds, role.env))
	}

	ctx.JSON(201, gin.H{"code": 201, "data": models.Role{}.FindByIdWithPermissions(r.ID, role.env)})
}

func (role *RoleController) Show(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("role"))
	if err != nil {
		log.Panicln("RoleController Show err: ", err)
		return
	}

	r := models.Role{}.FindByIdWithPermissions(uint(id), role.env)
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
		log.Panicln("RoleController Show err: ", err)
		return
	}
	if err := ctx.ShouldBind(&input); err != nil {
		exception.ValidateException(ctx, err, role.env)
		log.Println("RoleController Update err: ", err)
		return
	}

	r := models.Role{}.FindById(uint(id), role.env)
	if r == nil {
		exception.ModelNotFound(ctx, "role")
		return
	}

	if input.Name != "" {
		hasRole := models.Role{}.FindByName(input.Name, role.env)

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
		log.Println(input)
		e := role.env.GetDB().Model(r).Update("name", input.Name)
		log.Println(e.Error)
	}

	log.Println("input.PermissionIds", input.PermissionIds)
	if input.PermissionIds != nil {
		ps := models.Permission{}.FindByIds(input.PermissionIds, role.env)
		role.env.GetDB().Model(r).Association("Permissions").Clear()
		role.env.GetDB().Model(r).Association("Permissions").Append(toInterfaceSlice(ps)...)
	}

	ctx.JSON(200, gin.H{"code": 200, "data": models.Role{}.FindByIdWithPermissions(r.ID, role.env)})
}

func (role *RoleController) Destroy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("role"))
	if err != nil {
		log.Panicln("RoleController Destroy err: ", err)
		return
	}

	r := models.Role{}.FindById(uint(id), role.env)
	if r == nil {
		exception.ModelNotFound(ctx, "role")
		return
	}

	role.env.GetDB().Delete(r)
	ctx.JSON(204, nil)
}

func toInterfaceSlice(slice interface{}) []interface{} {
	permissions := slice.([]*models.Permission)
	newS := make([]interface{}, len(permissions))
	for i, v := range permissions {
		newS[i] = v
	}

	return newS
}
