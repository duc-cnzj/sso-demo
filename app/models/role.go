package models

import (
	"log"
	"sso/config/env"
	"time"
)

type Role struct {
	ID uint `gorm:"primary_key" json:"id"`

	Name string `gorm:"type:varchar(100);unique_index;not null;" json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permission;"`
}

func (Role) FindById(id uint, env *env.Env) *Role {
	r := &Role{}
	err := env.GetDB().Where("id = ?", id).First(r)
	if err.Error != nil {
		log.Println("findById", err)
		return nil
	}

	return r
}

func (r Role) FindByName(name string, env *env.Env)*Role {
	role := &Role{}
	err := env.GetDB().Where("name = ?", name).First(role)
	if err.Error != nil {
		log.Println("FindByName", err)
		return nil
	}
	log.Println(role)

	return role
}

func (r Role) FindByIdWithPermissions(id uint, env *env.Env)*Role {
	role := &Role{}
	err := env.GetDB().
		Preloads("Permissions").
		First(role, "id = ?", id)
	if err.Error != nil {
		log.Println("FindByIdWithPermissions", err)
		return nil
	}
	log.Println(role)

	return role
}
