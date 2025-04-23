package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
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

type UserRoleResponse struct {
	IDs    RoleIDs `json:"id"`
	RoleID string  `json:"role_id"`
}

func GetRoleID(db *gorm.DB, rID RoleIDs) (*UserRoleID, error) {
	var uRole = &UserRoleID{}
	query := `id = ?`

	r := strconv.Itoa(int(rID))
	preload := postgres.Preload(db, &UserRoleID{}, `Role`)
	err := postgres.SelectSingleRecord(preload, query, &UserRoleID{}, uRole, r)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return uRole, nil
}

func GetAllRoleIDs(db *gorm.DB) ([]UserRoleResponse, error) {
	roles := []UserRoleResponse{}
	err := postgres.SelectAllRecords(db, "", "", UserRoleID{}, &roles)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("no roles found")
		}
		return nil, err
	}
	return roles, nil
}

func GetAllRoles(db *gorm.DB) ([]Role, error) {
	roles := []Role{}
	err := postgres.SelectAllRecords(db, "", "", Role{}, &roles)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("no roles found")
		}
		return nil, err
	}
	return roles, nil
}
