package models

import (
	"time"
)

type Users struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"`
	RoleID    string    `json:"-" gorm:"not null"`
	Role      Role      `json:"role_id" gorm:"not null;foreignKey:RoleID;references:ID"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
