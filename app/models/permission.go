package models

import (
	"github.com/rs/zerolog/log"
	"sso/config/env"
	"time"
)

type Permission struct {
	ID uint `gorm:"primary_key" json:"id"`

	Name    string `gorm:"type:varchar(100);unique_index;not null;" json:"name"`
	Project string `json:"project"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Roles []Role `gorm:"many2many:role_permission;" json:"roles"`
	Users []User `gorm:"many2many:user_permission;" json:"users"`
}

func (r Permission) FindByIds(ids []uint, env *env.Env) []*Permission {
	var permissions []*Permission
	err := env.GetDB().Where("id in (?)", ids).Find(&permissions)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindByIds")
		return nil
	}

	return permissions
}

func (r Permission) FindById(id uint, env *env.Env) *Permission {
	var permission = &Permission{}
	err := env.GetDB().Where("id = ?", id).First(&permission)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindById")
		return nil
	}

	return permission
}

func (r Permission) FindByName(name string, env *env.Env) *Permission {
	var permission = &Permission{}
	err := env.GetDB().Where("name = ?", name).First(&permission)
	if err.Error != nil {
		log.Debug().Err(err.Error).Msg("Permission FindByName")

		return nil
	}

	return permission
}
