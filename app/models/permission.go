package models

import (
	"log"
	"sso/config/env"
	"time"
)

type Permission struct {
	ID uint `gorm:"primary_key" json:"id"`

	Name string `gorm:"type:varchar(100);unique_index;not null;" json:"name"`
	Project string `json:"project"`

	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         *time.Time     `sql:"index" json:"deleted_at"`

	Roles []Role `gorm:"many2many:role_permission;"`
}


func (r Permission) FindByIds(ids []uint, env *env.Env)[]*Permission {
	var permissions []*Permission
	err := env.GetDB().Where("id in (?)", ids).Find(&permissions)
	if err.Error != nil {
		log.Println("Permission FindByIds", err)
		return nil
	}

	return permissions
}