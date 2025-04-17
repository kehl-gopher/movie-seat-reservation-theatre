package auth

import (
	"errors"
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"gorm.io/gorm"
)

type UserAuthRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func ValidateUserAuthRequest(email, firstName, lastName, password string) (string, string, string, string, error) {
	v := utility.NewValidationError()

	firstName = v.ValidateName(firstName, "first_name")
	password = v.ValidatePassword(password)
	email = v.ValidateEmailAddress(email)
	lastName = v.ValidateName(lastName, "last_name")

	if v.CheckError() {
		return "", "", "", "", v
	}
	return firstName, email, lastName, password, nil
}

func UserRequestService(user UserAuthRequest, db *repository.Database, rIDs models.RoleIDs, expires_in int64, secret_key []byte) (int, *models.UserResponse, error) {
	firstName, email, lastName, password, err := ValidateUserAuthRequest(user.Email, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return http.StatusUnprocessableEntity, nil, err
	}
	password, err = utility.CreatePasswordHash(password)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	uRoleID, err := models.GetRoleID(db.Pdb.DB, rIDs)

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return http.StatusNotFound, nil, errors.New("role not found")
		} else if err.Error() == "validation error" {
			return http.StatusUnprocessableEntity, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}
	u := models.Users{
		ID:        utility.GenerateUUID(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		IsActive:  true,
		RoleID:    uRoleID.RoleID,
		Role:      uRoleID.Role,
	}

	userResponse, err := u.CreateUser(db, rIDs, expires_in, secret_key)
	if err != nil {
		if err.Error() == "validation error" {
			return http.StatusUnprocessableEntity, nil, err
		} else if errors.Is(gorm.ErrRecordNotFound, err) {
			return http.StatusNotFound, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusCreated, userResponse, nil
}

func AuthenticateUser(db *repository.Database, email, password string, secret_key []byte, exp_in int64, rID models.RoleIDs) (int, *models.UserResponse, error) {
	user := models.Users{
		Email:    email,
		Password: password,
	}
	u, err := user.GetUserByEmail(db)

	if err != nil || u == nil {
		if errors.Is(postgres.ErrNoRecordFound, err) || u == nil {
			return http.StatusBadRequest, nil, errors.New("invalid email or password")
		}
		return http.StatusInternalServerError, nil, err

	}
	checkPasswordHash, err := utility.VerifyPasswordHash(password, u.Password)
	if err != nil || !checkPasswordHash {
		if !checkPasswordHash {
			return http.StatusBadRequest, nil, errors.New("invalid email or password")
		}
		return http.StatusInternalServerError, nil, err
	}

	uRep, err := u.CreateUserSignInTokekn(db, secret_key, exp_in, rID)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, uRep, nil
}
