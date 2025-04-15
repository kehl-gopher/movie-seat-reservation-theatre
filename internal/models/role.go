package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleName string

type DefaultRole struct {
	User  RoleName
	Admin RoleName
}

var (
	User  RoleName = "user"
	Admin RoleName = "admin"
)

type Role struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// func (r *RoleName) GetRole() RoleName{
	
// }


