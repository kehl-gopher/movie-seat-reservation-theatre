package movies

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

const FilePath = "movies/images/"

type MovieReq struct {
	Title         string      `json:"title"`
	Synopsis      *string     `json:"synopsis"`
	GenreID       []string    `json:"genre_ids"`
	ReleasDate    models.Date `json:"release_date"`
	DurationInMin uint8       `json:"duration"`
	PosterImage   string      `json:"poster_image"`
	BackdropImage string      `json:"backdrop_image"`
}

// validate image base64 input
func ValidateMovieEntry(movie *MovieReq) (posterImage []byte, backDropImage []byte, ext_poster string, ext_backdrop string, e error) {
	v := &utility.ValidationError{}

	movie.Title = v.ValidateMovieName(movie.Title, "title")
	*movie.Synopsis = v.ValidateMovieSynopsis(*movie.Synopsis)
	movie.GenreID = v.ValidateMovieGenreID(movie.GenreID)
	movie.DurationInMin = v.ValidateMovieDuration(movie.DurationInMin)
	posterImage, ext_poster = v.ValidateImage(movie.PosterImage, "poster_image")
	backDropImage, ext_backdrop = v.ValidateImage(movie.PosterImage, "backdrop_image")

	if v.CheckError() {
		e = v
		return nil, nil, "", "", e
	}
	return posterImage, backDropImage, ext_poster, ext_backdrop, nil
}

func CreateMovies(db *repository.Database, movie *MovieReq, config *env.Config) (int, error) {

	posterImage, backdropImage, extPoster, extBackdrop, err := ValidateMovieEntry(movie)

	if err != nil {
		if err.Error() == "validation error" {
			return http.StatusUnprocessableEntity, err
		}
		return http.StatusInternalServerError, err
	}

	genres, err := models.GetGenresByID(db.Pdb.DB, movie.GenreID...)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRecordFound) {
			return http.StatusNotFound, errors.New("genre not found")
		}
		return http.StatusInternalServerError, err

	}

	m := &models.Movie{
		Title:       movie.Title,
		Synopsis:    *movie.Synopsis,
		ReleaseDate: movie.ReleasDate,
		Duration:    movie.DurationInMin,
		Genres:      genres,
	}

	pObjName := movie.Title + "_poster" + "." + extPoster
	bObjName := movie.Title + "_backdrop" + "." + extBackdrop
	err = m.CreateMovie(db, FilePath, config.MINIO_BUCKET, pObjName, bObjName, extPoster, extBackdrop, posterImage, backdropImage)

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}
