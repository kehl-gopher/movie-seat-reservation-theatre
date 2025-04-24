package movies

import (
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
)

func GetAllGenres(db *repository.Database) ([]models.Genre, int, error) {

	genre := &models.Genre{}
	genres, err := genre.GetAllGenres(db)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return genres, http.StatusOK, nil
}

func CreateGenres(db *repository.Database) {

}
