package permissioncontroller

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
	"strings"
)

type PermissionController struct {
	env *env.Env
	*api.AllRepo
}

type Uri struct {
	Permission uint `uri:"permission" binding:"required"`
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
	Sort    string `form:"sort"`

	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func NewPermissionController(env *env.Env) *PermissionController {
	return &PermissionController{env: env, AllRepo: api.NewAllRepo(env)}
}

func (p *PermissionController) Index(c *gin.Context) {
	var (
		query QueryInput
		count int
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		exception.ValidateException(c, err, p.env)
		return
	}
	log.Info().Interface("query", query).Msg("PermissionController.Index.query")

	var permissions []models.Permission
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

	q := p.env.GetDB().Model(&permissions)
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
		Order("id " + query.Sort).
		Find(&permissions)
	if len(permissions) < query.PageSize {
		count = query.PageSize*(query.Page-1) + len(permissions)
	} else {
		countQuery := p.env.GetDB().Model(&models.Permission{})
		if query.Name != "" {
			countQuery = countQuery.Where("name like ?", "%"+query.Name+"%")
		}
		if query.Project != "" {
			countQuery = countQuery.Where("project like ?", "%"+query.Project+"%")
		}
		countQuery.Count(&count)
	}

	c.JSON(200, gin.H{"code": 200, "data": permissions, "page": query.Page, "page_size": query.PageSize, "total": count})
}

func (p *PermissionController) Store(c *gin.Context) {
	var input StoreInput
	if err := c.ShouldBind(&input); err != nil {
		exception.ValidateException(c, err, p.env)
		log.Debug().Err(err).Msg("PermissionController Store err: ")
		return
	}

	permission, _ := p.PermRepo.FindByName(input.Name)
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

	pnew := &models.Permission{
		Name:    input.Name,
		Project: input.Project,
	}
	if err := p.PermRepo.Create(pnew); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	c.JSON(201, gin.H{"code": 201, "data": pnew})
}

func (p *PermissionController) Show(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		return
	}

	r, _ := p.PermRepo.FindById(uri.Permission)
	if r == nil {
		exception.ModelNotFound(c, "permission")
		return
	}

	c.JSON(200, gin.H{"data": r})
}

func (p *PermissionController) Update(c *gin.Context) {
	var (
		input UpdateInput
		uri   Uri
	)
	if err := c.ShouldBindUri(&uri); err != nil {
		return
	}
	if err := c.ShouldBind(&input); err != nil {
		exception.ValidateException(c, err, p.env)
		log.Debug().Err(err).Msg("PermissionController Update err: ")
		return
	}
	permission, _ := p.PermRepo.FindById(uri.Permission)
	if permission == nil {
		exception.ModelNotFound(c, "Permission")
		return
	}
	hasPermission, _ := p.PermRepo.FindByName(input.Name)
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

	e := p.env.GetDB().Model(permission).Updates(map[string]interface{}{
		"name": input.Name, "project": input.Project,
	})
	if e.Error != nil {
		log.Panic().Msg(e.Error.Error())
	}
	c.JSON(200, gin.H{"code": 200, "data": permission})
}

func (p *PermissionController) Destroy(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithError(500, err)
		return
	}

	r, _ := p.PermRepo.FindById(uri.Permission)

	if r == nil {
		exception.ModelNotFound(c, "Permission")
		return
	}

	p.env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Delete(r).Error != nil {
			return tx.Delete(r).Error
		}
		if err := tx.Model(r).Association("Roles").Clear().Error; err != nil {
			return err
		}

		return nil
	})

	c.JSON(204, nil)
}

func (p *PermissionController) GetByGroups(c *gin.Context) {
	var permissions []models.Permission

	type Item struct {
		Name string `json:"name"`
		Id   uint   `json:"id"`
	}
	var groups = make(map[string][]Item)

	p.env.GetDB().
		Order("id DESC").
		Find(&permissions)
	for _, permission := range permissions {
		if items, ok := groups[permission.Project]; ok {
			groups[permission.Project] = append(items, Item{
				Name: permission.Name,
				Id:   permission.ID,
			})
		} else {
			groups[permission.Project] = []Item{
				{
					Name: permission.Name,
					Id:   permission.ID,
				}}
		}

	}

	c.JSON(200, gin.H{"data": groups})
}

func (p *PermissionController) GetPermissionProjects(c *gin.Context) {
	var res []models.Permission
	if err := p.env.GetDB().Model(&models.Permission{}).Select([]string{"distinct project"}).Find(&res).Error; err != nil {
		log.Panic().Err(err).Msg("PermissionController.GetPermissionProjects")
	}

	var items []string
	for _, value := range res {
		items = append(items, value.Project)
	}

	c.JSON(200, gin.H{"data": items})
}
