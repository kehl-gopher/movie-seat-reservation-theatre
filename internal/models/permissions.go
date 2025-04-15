package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	PermissionList PermissionList `gorm:"type:jsonb;not null"`
	Role           Role           `json:"role_id" gorm:"not null;unique"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type PermissionList struct {
	CanCreateMovies  bool `json:"can_create_movies"`
	CanDeleteMovie   bool `json:"can_delete_movies"`
	CanUpdateMovie   bool `json:"can_update_movies"`
	CanGetMovies     bool `json:"can_get_movies"`
	CanDeleteUsers   bool `json:"can_delete_users"`
	CanGetUsers      bool `json:"can_get_users"`
	CanBanUsers      bool `json:"can_ban_users"`
	CanCreateUsers   bool `json:"can_create_users"`
	CanBookSeats     bool `json:"can_book_seats"`
	CanCreateSeats   bool `json:"can_create_seats"`
	CanRemoveSeats   bool `json:"can_delete_seats"`
	CanUpdateSeats   bool `json:"can_update_seats"`
	CanCancelBooking bool `json:"can_cancel_booking"`
	CanGetBookings   bool `json:"can_get_bookings"`
	CanGetRoles      bool `json:"can_get_roles"`
	CanCreateRoles   bool `json:"can_create_roles"`
	CanUpdateRoles   bool `json:"can_update_roles"`
	CanDeleteRoles   bool `json:"can_delete_roles"`
}
