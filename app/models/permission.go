package models

import (
	"log"
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
		log.Println("Permission FindByIds", err)
		return nil
	}

	return permissions
}

func (r Permission) FindById(id uint, env *env.Env) *Permission {
	var permission = &Permission{}
	err := env.GetDB().Where("id = ?", id).First(&permission)
	if err.Error != nil {
		log.Println("Permission FindById", err)
		return nil
	}

	return permission
}

func (r Permission) FindByName(name string, env *env.Env) *Permission {
	var permission = &Permission{}
	err := env.GetDB().Where("name = ?", name).First(&permission)
	if err.Error != nil {
		log.Println("Permission FindByName", err)
		return nil
	}

	return permission
}
