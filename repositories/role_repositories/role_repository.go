package role_repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
	"sso/repositories/permission_repository"
)

type RoleRepository struct {
	env *env.Env
}

func NewRoleRepository(env *env.Env) *RoleRepository {
	return &RoleRepository{env: env}
}

func (repo *RoleRepository) FindByIds(ids []uint) ([]*models.Role, error) {
	var roles []*models.Role
	if err := repo.env.GetDB().Where("id in (?)", ids).Find(&roles).Error; err != nil {
		log.Debug().Err(err).Msg("FindByIds")
		return nil, err
	}

	return roles, nil
}

func (repo *RoleRepository) FindById(id uint) (*models.Role, error) {
	r := &models.Role{}
	if err := repo.env.GetDB().Where("id = ?", id).First(r).Error; err != nil {
		log.Debug().Err(err).Msg("findById")
		return nil, err
	}

	return r, nil
}

func (repo *RoleRepository) FindByName(name string) (*models.Role, error) {
	role := &models.Role{}

	if err := repo.env.GetDB().Where("name = ?", name).First(role).Error; err != nil {
		log.Debug().Err(err).Msg("FindByName")

		return nil, err
	}
	log.Debug().Interface("role", role).Msg("FindByName")

	return role, nil
}

func (repo *RoleRepository) FindByIdWithPermissions(id uint) (*models.Role, error) {
	role := &models.Role{}

	if err := repo.env.GetDB().
		Preload("Permissions").
		First(role, "id = ?", id).Error; err != nil {
		log.Debug().Err(err).Msg("FindByIdWithPermissions")
		return nil, err
	}
	log.Debug().Interface("role", role).Msg("FindByName")

	return role, nil
}

func (repo *RoleRepository) Create(r *models.Role) error {
	if err := repo.env.GetDB().Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RoleRepository) SyncPermissions(role *models.Role, ids []uint, db *gorm.DB) error {
	var conn *gorm.DB

	if db != nil {
		conn = db
	} else {
		conn = repo.env.GetDB()
	}
	if ids != nil {
		permRepo := permission_repository.NewPermissionRepository(repo.env)
		permissionByIds, _ := permRepo.FindByIds(ids)
		if err := conn.Model(role).Association("Permissions").Clear().Error; err != nil {
			return err
		}

		if err := conn.Model(role).Association("Permissions").Append(toInterfaceSlice(permissionByIds)...).Error; err != nil {
			return err
		}
	}

	return nil
}

func toInterfaceSlice(slice interface{}) []interface{} {
	permissions := slice.([]*models.Permission)
	newS := make([]interface{}, len(permissions))
	for i, v := range permissions {
		newS[i] = v
	}

	return newS
}
