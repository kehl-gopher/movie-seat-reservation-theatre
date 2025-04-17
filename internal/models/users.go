package models

import (
	"fmt"
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

type AccessToken struct {
	ID          string    `json:"-" gorm:"primaryKey"`
	Token       string    `json:"token" gorm:"not null"`
	Expiry      time.Time `json:"expires_at" gorm:"not null"`
	BlackListed bool      `json:"-" gorm:"default:false"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
	UserID      string    `json:"-" gorm:"not null"`
	User        Users     `json:"-" gorm:"not null;foreignKey:UserID;references:ID"`
}

type UserResponse struct {
	ID          string      `json:"id"`
	Email       string      `json:"email"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Role        uint8       `json:"role"`
	IsActive    bool        `json:"is_active"`
	AccessToken AccessToken `json:"access_token"`
}

func (u *Users) CreateUser(db *repository.Database, uRoleID RoleIDs, exp_in int64, secret_key []byte) (*UserResponse, error) {
	if userExist, err := postgres.CheckExists(db.Pdb.DB, `email = ?`, &Users{}, u.Email); err == nil && userExist {
		v := utility.NewValidationError()
		v.AddValidationError("email", "already exists")
		return nil, v
	} else if err != nil {
		return nil, err
	}

	err := postgres.Create(db.Pdb.DB, u)

	if err != nil {
		return nil, err
	}

	accessToken, err := u.CreateAccessToken(db.Pdb.DB, secret_key, exp_in, uRoleID)
	if err != nil {
		return nil, err
	}
	uResponse := &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Role:        uint8(uRoleID),
		IsActive:    u.IsActive,
		AccessToken: accessToken,
	}
	return uResponse, nil
}

func (u *Users) CreateAccessToken(db *gorm.DB, secret_key []byte, expires_in int64, rIDs RoleIDs) (AccessToken, error) {
	claims := utility.AccessTokenClaim{
		UserId:    u.ID,
		Role:      uint8(rIDs),
		SecretKey: secret_key,
		ExpiresAt: expires_in,
	}
	token, err := claims.CreateNewToken()
	if err != nil {
		return AccessToken{}, err
	}
	fmt.Println(u.ID)
	accessToken := AccessToken{
		ID:          utility.GenerateUUID(),
		BlackListed: false,
		Token:       token,
		UserID:      u.ID,
		User:        *u,
		Expiry:      time.Now().Add(time.Duration(expires_in) * time.Second),
	}

	// check if user token exists if its does delete and create a new one
	query := "user_id = ?"
	err = postgres.DeleteSingleRecord(db, query, &AccessToken{}, u.ID)
	if err != nil {
		return AccessToken{}, err
	}
	err = postgres.Create(db, &accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}

func (u *Users) GetUserByEmail(db *repository.Database) (*Users, error) {
	var user = &Users{}

	query := `email = ?`
	err := postgres.SelectSingleRecord(db.Pdb.DB, query, &Users{}, user, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) GetUserById(db *repository.Database) (*Users, error) {
	var user = &Users{}
	err := postgres.SelectById(db.Pdb.DB, u.ID, &Users{}, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) CreateUserSignInTokekn(db *repository.Database, secret_key []byte, exp_in int64, uRoleID RoleIDs) (*UserResponse, error) {
	accessToken, err := u.CreateAccessToken(db.Pdb.DB, secret_key, exp_in, uRoleID)
	if err != nil {
		return nil, err
	}
	uResponse := &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Role:        uint8(uRoleID),
		IsActive:    u.IsActive,
		AccessToken: accessToken,
	}
	return uResponse, nil
}
