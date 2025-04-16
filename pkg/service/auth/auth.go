package auth

import (
	"fmt"
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

type UserAuthRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func ValidateUserAuthRequest(email, firstName, lastName, password string) (string, string, string, string, *utility.ValidationError) {
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

func UserRequestService(user UserAuthRequest, db repository.Database) (int, error) {

	firstName, email, lastName, password, err := ValidateUserAuthRequest(user.Email, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	u := models.Users{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		IsActive:  true,
	}
	fmt.Printf("%+v", u)
	return 0, nil
}
