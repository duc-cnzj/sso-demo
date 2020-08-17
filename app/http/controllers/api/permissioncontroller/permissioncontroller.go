package permissioncontroller

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

type PermissionController struct {
	env *env.Env
}

type StoreInput struct {
	Name    string `form:"name" json:"name"`
	Project string `form:"project" json:"project"`
}

type UpdateInput struct {
	Name    string `form:"name" json:"name"`
	Project string `form:"project" json:"project"`
}

type QueryInput struct {
	Name    string `form:"name"`
	Project string `form:"project"`

	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func NewPermissionController(env *env.Env) *PermissionController {
	return &PermissionController{env: env}
}

func (p *PermissionController) Index(c *gin.Context) {
	var query QueryInput
	if err := c.ShouldBindQuery(&query); err != nil {
		exception.ValidateException(c, err, p.env)
		return
	}
	log.Println(query)

	var roles []models.Permission
	if query.PageSize <= 0 {
		query.PageSize = 15
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	q := p.env.GetDB().Model(&roles)
	if query.Name != "" {
		q = q.Where("name like ?", "%"+query.Name+"%")
	}
	if query.Project != "" {
		q = q.Where("project like ?", "%"+query.Project+"%")
	}

	offset := int(math.Max(float64((query.Page-1)*query.PageSize), 0))
	q.
		Offset(offset).
		Limit(query.PageSize).
		Order("id DESC").
		Find(&roles)

	c.JSON(200, gin.H{"code": 200, "data": roles})
}

func (p *PermissionController) Store(c *gin.Context) {
	var input StoreInput
	if err := c.ShouldBind(input); err != nil {
		exception.ValidateException(c, err, p.env)
		log.Println("PermissionController Store err: ", err)
		return
	}

	permission := models.Permission{}.FindByName(input.Name, p.env)
	if permission != nil {
		var errors = form.ValidateErrors{
			form.ValidateError{
				Field: "name",
				Msg:   "permission exists",
			},
		}
		exception.ValidateException(c, errors, p.env)
		return
	}

	pnew := models.Permission{
		Name:    input.Name,
		Project: input.Project,
	}
	if newP := p.env.GetDB().Create(pnew); newP.Error != nil {
		log.Panicln(newP.Error.Error())
	}

	c.JSON(201, gin.H{"code": 201, "data": pnew})
}

func (p *PermissionController) Show(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("permission"))
	if err != nil {
		log.Panicln("PermissionController Show err: ", err)
		return
	}

	r := models.Permission{}.FindById(uint(id), p.env)
	if r == nil {
		exception.ModelNotFound(c, "permission")
		return
	}

	c.JSON(200, gin.H{"data": r})
}

func (p *PermissionController) Update(c *gin.Context) {
	var input UpdateInput
	id, err := strconv.Atoi(c.Param("permission"))
	if err != nil {
		log.Panicln("PermissionController Show err: ", err)
		return
	}
	if err := c.ShouldBind(&input); err != nil {
		exception.ValidateException(c, err, p.env)
		log.Println("PermissionController Update err: ", err)
		return
	}
	permission := models.Permission{}.FindById(uint(id), p.env)
	if permission == nil {
		exception.ModelNotFound(c, "Permission")
		return
	}
	hasPermission := models.Role{}.FindByName(input.Name, p.env)
	if hasPermission != nil && hasPermission.ID != permission.ID {
		var errors = form.ValidateErrors{
			form.ValidateError{
				Field: "name",
				Msg:   "name exists",
			},
		}
		exception.ValidateException(c, errors, p.env)
		return
	}

	e := p.env.GetDB().Model(permission).Update("name", input.Name, "project", input.Project)
	if e.Error != nil {
		log.Panicln(e.Error.Error())
	}
	c.JSON(200, gin.H{"code": 200, "data": permission})
}

func (p *PermissionController) Destroy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("permission"))
	if err != nil {
		log.Panicln("PermissionController Destroy err: ", err)
		return
	}

	r := models.Permission{}.FindById(uint(id), p.env)
	if r == nil {
		exception.ModelNotFound(c, "Permission")
		return
	}

	p.env.GetDB().Delete(r)
	c.JSON(204, nil)
}
