package role_repositories

import (
	"github.com/rs/zerolog/log"
	"sso/app/models"
	"sso/config/env"
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
