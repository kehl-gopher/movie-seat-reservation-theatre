package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleName string
type RoleIDs uint8

type DefaultRole struct {
	User  RoleIDs
	Admin RoleIDs
}

var (
	User      RoleName = "user"
	Admin     RoleName = "admin"
	Anonymous RoleName = "anonymous"
)

type Role struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserRoleID struct {
	ID     uint8  `json:"id" gorm:"primaryKey"`
	RoleID string `json:"role_id" gorm:"not null"`
	Role   Role   `json:"role" gorm:"not null;foreignKey:RoleID;references:ID"`
}

func (r RoleIDs) String() RoleName {
	switch r {
	case 1:
		return Admin
	case 2:
		return User
	default:
		return Anonymous
	}
}

func (r *Role) GetUserRoleID(db *gorm.DB) {

}
