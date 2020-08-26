package models

import (
	"github.com/rs/zerolog/log"
	"sso/config/env"
	"time"
)

type Role struct {
	ID uint `gorm:"primary_key" json:"id"`

	Name string `gorm:"type:varchar(100);unique_index;not null;" json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permission;" json:"permissions"`
	Users       []User       `gorm:"many2many:user_role;" json:"users"`
}

func (Role) FindByIds(ids []uint, env *env.Env) []*Role {
	var roles []*Role
	if err := env.GetDB().Where("id in (?)", ids).Find(&roles).Error; err != nil {
		log.Debug().Err(err).Msg("FindByIds")
		return nil
	}

	return roles
}

func (Role) FindById(id uint, env *env.Env) *Role {
	r := &Role{}
	if err := env.GetDB().Where("id = ?", id).First(r).Error; err != nil {
		log.Debug().Err(err).Msg("findById")
		return nil
	}

	return r
}

func (r Role) FindByName(name string, env *env.Env) *Role {
	role := &Role{}

	if err := env.GetDB().Where("name = ?", name).First(role).Error; err != nil {
		log.Debug().Err(err).Msg("FindByName")

		return nil
	}
	log.Debug().Interface("role", role).Msg("FindByName")

	return role
}

func (r Role) FindByIdWithPermissions(id uint, env *env.Env) *Role {
	role := &Role{}

	if err := env.GetDB().
		Preload("Permissions").
		First(role, "id = ?", id).Error; err != nil {
		log.Debug().Err(err).Msg("FindByIdWithPermissions")
		return nil
	}
	log.Debug().Interface("role", role).Msg("FindByName")

	return role
}
