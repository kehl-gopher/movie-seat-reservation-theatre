package seathalls

import (
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

func ValidateHallInputs(hallName string, rows, numberOfSeats int) error {
	v := utility.NewValidationError()
	if hallName == "" {
		v.AddValidationError("hall_name", "field is required")
	} else {
		v.ValidateHallName(hallName)
	}
	if rows <= 0 {
		v.AddValidationError("number_of_rows", "field is required")
	}
	if numberOfSeats <= 0 {
		v.AddValidationError("number_of_seats", "field is required")
	}

	if v.CheckError() {
		return v
	}
	return nil
}
func CreateHallSeat(db *repository.Database, config *env.Config, hallName string, rows, numberOfSeats int) (*models.Halls, int, error) {

	err := ValidateHallInputs(hallName, rows, numberOfSeats)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	halls := models.Halls{Name: hallName}
	data, statusCode, err := halls.CreateHallSeat(db, config, rows, numberOfSeats)
	if err != nil {
		return nil, statusCode, err
	}
	return data, statusCode, nil
}

func ValidateHallUpdateInputs(hallName string, rows, numberOfSeats int) error {
	v := utility.NewValidationError()

	if rows <= 0 {
		v.AddValidationError("number_of_rows", "field is required")
	}
	if numberOfSeats <= 0 {
		v.AddValidationError("number_of_seats", "field is required")
	}

	if v.CheckError() {
		return v
	}
	return nil
}

// TODO implement user perform update on seats
func UpdateHallSeat(db *repository.Database, config *env.Config, hallID string, hallName *string, rows, numberOfSeats *int) (*models.Halls, int, error) {

	// get hall
	return nil, 0, nil
}
