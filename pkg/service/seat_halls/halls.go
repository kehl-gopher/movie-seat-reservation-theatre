package seathalls

import (
	"errors"
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
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

// TBF there's no reason for me to be returning a reference of the object
// it's just plain out laziness in my CodeBase of having to write the whole code
// sheeeessh lazy fuck ain't it
func GetAllHalls(db *repository.Database, config *env.Config) ([]models.Halls, int, error) {
	h := models.Halls{}

	halls, err := h.GetAllHalls(db)
	if err != nil {
		if errors.Is(postgres.ErrNoRecordFound, err) {
			return nil, http.StatusNotFound, errors.New("no Halls found")
		}
		return nil, http.StatusInternalServerError, err
	}
	return halls, http.StatusCreated, nil
}

func GetHallDetails(db *repository.Database, config *env.Config, hallId string) (models.Halls, int, error) {

	h := models.Halls{ID: hallId}

	hall, err := h.GetAllDetails(db)
	if err != nil {
		if errors.Is(postgres.ErrNoRecordFound, err) {
			return models.Halls{}, http.StatusNotFound, err
		}
	}
	return hall, http.StatusOK, nil
}

// sheeesh so dumb he could not figure out how to handle... validation logic for simple hall updates... FUCK...
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
