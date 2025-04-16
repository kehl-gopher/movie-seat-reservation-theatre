package models

import (
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"gorm.io/gorm"
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

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utility.GenerateUUID()
	return
}

func (u *Users) CreateUser(db *repository.Database) error {
	if err := postgres.Create(db.Pdb.DB, u); err != nil {
		return err
	}
	return nil
}
