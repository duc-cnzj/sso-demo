package permission_repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
)

type PermissionRepositoryImp interface {
	FindByIds([]uint) ([]*models.Permission, error)
	FindById(uint) (*models.Permission, error)
	FindByName(string) (*models.Permission, error)
	Create(*models.Permission) error
}

type PermissionRepository struct {
	env *env.Env
}

func NewPermissionRepository(env *env.Env) *PermissionRepository {
	return &PermissionRepository{env: env}
}

func (repo *PermissionRepository) FindByIds(ids []uint) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := repo.env.GetDB().Where("id in (?)", ids).Find(&permissions)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindByIds")
		return nil, err.Error
	}

	return permissions, nil
}

func (repo *PermissionRepository) FindById(id uint) (*models.Permission, error) {
	var permission = &models.Permission{}
	err := repo.env.GetDB().Where("id = ?", id).First(&permission)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindById")
		return nil, err.Error
	}

	return permission, nil
}

func (repo *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission = &models.Permission{}
	err := repo.env.GetDB().Where("name = ?", name).First(&permission)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindByName")

		return nil, err.Error
	}

	return permission, nil
}

func (repo *PermissionRepository) Create(permission *models.Permission) error {
	var p = &models.Permission{}
	if err := repo.env.GetDB().First(p, "name = ? AND project = ?", permission.Name, permission.Project).Error; err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		if err := repo.env.GetDB().Create(permission).Error; err != nil {
			return err
		}
	} else {
		*permission = *p
	}

	return nil
}