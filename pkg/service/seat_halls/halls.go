package seathalls

import (
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
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
func CreateHallSeat(db *repository.Database, env *env.Config, hallName string, rows, numberOfSeats int) (int, error) {

	err := ValidateHallInputs(hallName, rows, numberOfSeats)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}
	
	
	
	return 0, nil
}
