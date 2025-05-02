package movies

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

const FilePath = "movies/images/"

type MovieReq struct {
	Title         string       `json:"title"`
	Synopsis      *string      `json:"synopsis"`
	GenreID       []string     `json:"genre_ids"`
	ReleasDate    *models.Date `json:"release_date"`
	DurationInMin uint8        `json:"duration"`
	PosterImage   string       `json:"poster_image"`
	BackdropImage string       `json:"backdrop_image"`
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
	if movie.ReleasDate == nil {
		v.AddValidationError("release_date", "movie release date is required")
	}

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
		ReleaseDate: *movie.ReleasDate,
		Duration:    models.Duration(movie.DurationInMin),
		Genres:      genres,
	}

	pObjName := movie.Title + "_poster" + utility.GenerateUUID() + "." + extPoster
	bObjName := movie.Title + "_backdrop" + utility.GenerateUUID() + "." + extBackdrop
	err = m.CreateMovie(db, FilePath, config.MINIO_BUCKET, pObjName, bObjName, extPoster, extBackdrop, posterImage, backdropImage)

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func GetMovieByID(db *repository.Database, id string) (*models.Movie, int, error) {

	movie := &models.Movie{}

	mov, err := movie.GetDetailMovie(db, id)
	if err != nil {
		if err.Error() == "movie not found" {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	return mov, http.StatusOK, nil
}

func cleanParams(search, date string, genres []string) (searchTerm string, genre string, release_date string) {
	if search != "" {
		searchTerm = strings.TrimLeft(search, "\"")
	} else {
		searchTerm = ""
	}
	if date != "" {
		date, _ := time.Parse("2006-01-02", date)
		release_date = date.Format("2006-01-02")
	} else {
		release_date = ""
	}
	if genres != nil {
		for ind, gen := range genres {
			genres[ind] = strings.Trim(strings.Title(gen), "\"")
		}
		genre = strings.Join(genres, ",")
	} else {
		genre = ""
	}

	return searchTerm, genre, release_date
}
func GetAllMovies(db *repository.Database, offset, limit uint, config *env.Config, search string, date string, genre []string) ([]models.MovieResponse, int, postgres.PaginationResponse, error) {

	movie := models.Movie{}

	search, genres, datet := cleanParams(search, date, genre)
	mov, pag, err := movie.GetAllMoviesWithPagination(db, offset, limit, config, search, datet, genres)

	if err != nil {
		if err.Error() == "no movies found" {
			return nil, http.StatusNotFound, pag, err
		}
		return nil, http.StatusInternalServerError, pag, err
	}
	return mov, http.StatusOK, pag, nil
}

type UpdateMovieReq struct {
	Title         *string      `json:"title"`
	Synopsis      *string      `json:"synopsis"`
	GenreID       []string     `json:"genre_ids"`
	ReleasDate    *models.Date `json:"release_date"`
	DurationInMin *uint8       `json:"duration"`
	PosterImage   *string      `json:"poster_image"`
	BackdropImage *string      `json:"backdrop_image"`
}

func ValidateUpdateMovieEntry(movie *UpdateMovieReq) (posterImage []byte, backDropImage []byte, ext_poster string, ext_backdrop string, e error) {
	v := &utility.ValidationError{}

	if movie.Title != nil {
		*movie.Title = v.ValidateMovieName(*movie.Title, "title")
	}
	if movie.Synopsis != nil {
		*movie.Synopsis = v.ValidateMovieSynopsis(*movie.Synopsis)
	}
	if movie.GenreID != nil {
		movie.GenreID = v.ValidateMovieGenreID(movie.GenreID)
	}
	if movie.DurationInMin != nil {
		*movie.DurationInMin = v.ValidateMovieDuration(*movie.DurationInMin)
	}
	if movie.PosterImage != nil {
		posterImage, ext_poster = v.ValidateImage(*movie.PosterImage, "poster_image")
	}
	if movie.BackdropImage != nil {
		backDropImage, ext_backdrop = v.ValidateImage(*movie.BackdropImage, "backdrop_image")
	}

	if v.CheckError() {
		e = v
		return nil, nil, "", "", e
	}
	return posterImage, backDropImage, ext_poster, ext_backdrop, nil
}

func UpdateMovie(db *repository.Database, movie_id string, movie *UpdateMovieReq, config *env.Config) (int, error) {
	var pObjName, bObjName, oldOrNewTitle string
	m := &models.Movie{}
	mov, err := m.GetMovieByID(db, movie_id)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRecordFound) {
			return http.StatusNotFound, errors.New("movie not found")
		}
		return http.StatusInternalServerError, err
	}

	posterImage, backdropImage, extPoster, extBackdrop, err := ValidateUpdateMovieEntry(movie)

	if err != nil {
		if err.Error() == "validation error" {
			return http.StatusUnprocessableEntity, err
		}
		return http.StatusInternalServerError, err
	}

	if movie.Title != nil {
		oldOrNewTitle = mov.Title
		mov.Title = *movie.Title
	}
	if movie.Synopsis != nil {
		mov.Synopsis = *movie.Synopsis
	}
	if movie.GenreID != nil {
		genres, err := models.GetGenresByID(db.Pdb.DB, movie.GenreID...)
		if err != nil {
			if errors.Is(err, postgres.ErrNoRecordFound) {
				return http.StatusNotFound, errors.New("genre not found")
			}
			return http.StatusInternalServerError, err
		}
		mov.Genres = genres
	}
	if movie.DurationInMin != nil {
		mov.Duration = models.Duration(*movie.DurationInMin)
	}
	if movie.ReleasDate != nil {
		mov.ReleaseDate = *movie.ReleasDate
	}
	if movie.PosterImage != nil {
		uuID := strings.Split(strings.Split(mov.PosterPath, "_poster")[1], ".")[0] // extract old uuid from image path
		if movie.Title == nil {
			oldOrNewTitle = mov.Title
		}
		pObjName = oldOrNewTitle + "_poster" + uuID + "." + extPoster
	}
	if movie.BackdropImage != nil {
		uuID := strings.Split(strings.Split(mov.BackDropPath, "_backdrop")[1], ".")[0] // extract old uuid from image path
		if movie.Title == nil {
			oldOrNewTitle = mov.Title
		}
		bObjName = oldOrNewTitle + "_backdrop" + uuID + "." + extBackdrop
	}
	err = mov.UpdateMovie(db, FilePath, config.MINIO_BUCKET, pObjName, bObjName, extPoster, extBackdrop, posterImage, backdropImage)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func DeleteMovie(db *repository.Database, movie_id string) (int, error) {
	m := &models.Movie{}
	mov, err := m.GetMovieByID(db, movie_id)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRecordFound) {
			return http.StatusNotFound, errors.New("movie not found")
		}
		return http.StatusInternalServerError, err
	}

	err = mov.DeleteMovie(db)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
