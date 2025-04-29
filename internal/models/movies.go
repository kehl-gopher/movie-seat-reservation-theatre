package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/minio"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	mini "github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type Date time.Time // define custom date for release date
type Duration uint8

type Movie struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"type:text;not null;index"`
	Synopsis     string         `json:"synopsis" gorm:"type:text;index"`
	BackDropPath string         `json:"backdrop_path" gorm:"type:text"`
	PosterPath   string         `json:"poster_path" gorm:"type:text"`
	ReleaseDate  Date           `json:"release_date" gorm:"type:date;not null"`
	Duration     Duration       `json:"duration" gorm:"type:smallint;not null"`
	CreatedAt    time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Genres       []Genre        `json:"genres" gorm:"many2many:movie_genres;"`
	Shows        []Shows        `json:"shows" gorm:"foreignKey:MovieID;references:ID"`
}

type Genre struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:text;not null;unique;index"`
}

// type MovieDetailResponse struct {
// 	ID           string         `json:"id"`
// 	Title        string         `json:"title"`
// 	Synopsis     string         `json:"synopsis" `
// 	BackDropPath string         `json:"backdrop_path"`
// 	PosterPath   string         `json:"poster_path"`
// 	ReleaseDate  Date           `json:"release_date"`
// 	Duration     Duration       `json:"duration"`
// 	CreatedAt    time.Time      `json:"-"`
// 	UpdatedAt    time.Time      `json:"-"`
// 	DeletedAt    gorm.DeletedAt `json:"-"`
// 	Genres       []Genre        `json:"genres"`
// 	Shows        []Shows        `json:"shows"`
// }

type MovieResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	PosterPath  string   `json:"poster_path"`
	ReleaseDate Date     `json:"release_date"`
	Duration    Duration `json:"duration"`
	Url         string   `json:"url" gorm:"-"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) error {
	m.ID = utility.GenerateUUID()
	return nil
}

func UploadImageToMinio(min *mini.Client, imageByte []byte, filePath string, ext string, bucketName, objName string) (string, error) {

	contentType := "image/" + ext
	url, err := minio.UploadToMinio(min, filePath, bucketName, contentType, objName, imageByte)

	if err != nil {
		return "", err
	}
	return url, nil
}

func (m *Movie) CreateMovie(db *repository.Database, filePath string, bucketName string, obj1 string, obj2 string, ext1 string, ext2 string, imageBytes ...[]byte) error {
	// generate minio object url
	profilePath, err := UploadImageToMinio(db.Min, imageBytes[0], filePath, ext1, bucketName, obj1)
	if err != nil {
		return nil
	}
	backdropPath, err := UploadImageToMinio(db.Min, imageBytes[1], filePath, ext2, bucketName, obj2)
	if err != nil {
		return err
	}

	m.PosterPath = profilePath
	m.BackDropPath = backdropPath

	if err := postgres.Create(db.Pdb.DB, m); err != nil {
		return err
	}
	return nil
}

func GetGenresByID(db *gorm.DB, genreIDs ...string) ([]Genre, error) {
	var genre = []Genre{}
	query := `id IN ?`
	err := postgres.SelectMultipleRecord(db, query, &Genre{}, &genre, genreIDs)

	return genre, err
}

func (m *Genre) GetAllGenres(db *repository.Database) ([]Genre, error) {
	var genres []Genre
	err := postgres.SelectAllRecords(db.Pdb.DB, "", "name", &Genre{}, &genres)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("no genres found")
		}
		return nil, err
	}
	return genres, nil
}

// TODO: get all movie relations
func (m *Movie) GetDetailMovie(db *repository.Database, id string) (*Movie, error) {
	movie := &Movie{}
	preload := postgres.Preload(db.Pdb.DB, movie, `Genres`, `Shows`, `Shows.Hall`, `Shows.Hall.Seats`)

	err := postgres.SelectSingleRecord(preload, `id = ?`, m, movie, id)
	if err != nil {
		if errors.Is(postgres.ErrNoRecordFound, err) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return movie, nil
}

func (m *Movie) GetMovieByID(db *repository.Database, id string) (*Movie, error) {
	movie := &Movie{}
	err := postgres.SelectSingleRecord(db.Pdb.DB, `id = ?`, m, movie, id)
	if err != nil {
		if errors.Is(postgres.ErrNoRecordFound, err) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return movie, nil

}

// TODO: get all movies with paginations paginations
func (m *Movie) GetAllMoviesWithPagination(db *repository.Database, page, limit uint, config *env.Config) ([]MovieResponse, postgres.PaginationResponse, error) {
	var movies []MovieResponse

	pag, err := postgres.SelectAllRecordWithPagination(db.Pdb.DB, "", "desc", "created_at", &Movie{}, &movies, limit, page)

	if err != nil {
		return nil, pag, err
	}
	if len(movies) == 0 {
		return nil, pag, errors.New("no movies found")
	}

	for ind, movie := range movies {
		movies[ind] =
			MovieResponse{
				ID:          movie.ID,
				Title:       movie.Title,
				PosterPath:  movie.PosterPath,
				ReleaseDate: movie.ReleaseDate,
				Duration:    movie.Duration,
				Url:         fmt.Sprintf("%s%s/movie/%s", config.APP_URL, config.BASEURL, movie.ID),
			}
	}

	return movies, pag, nil
}

// TODO: update movies
func (m *Movie) UpdateMovie(db *repository.Database, filePath string, bucketName string, obj1 string, obj2 string, ext1 string, ext2 string, imageBytes ...[]byte) error {
	// movie := &Movie{}
	var err error
	var profilePath, backdropPath string
	if obj1 != "" {
		fmt.Println(obj1)
		profilePath, err = UploadImageToMinio(db.Min, imageBytes[0], filePath, ext1, bucketName, obj1)
		if err != nil {
			return nil
		}
	}
	if obj2 != "" {
		backdropPath, err = UploadImageToMinio(db.Min, imageBytes[1], filePath, ext2, bucketName, obj2)
		if err != nil {
			return err
		}

	}

	m.PosterPath = profilePath
	m.BackDropPath = backdropPath

	query := `id = ?`
	err = postgres.UpdateSingleRecord(db.Pdb.DB, query, m, m.ID)

	if err != nil {
		return err
	}

	genres := map[string]interface{}{
		"Genres": m.Genres,
	}
	err = postgres.UpdateRelationShipRecord(db.Pdb.DB, m, genres)

	if err != nil {
		return err
	}
	return nil
}

// TODO: delete movie
func (m *Movie) DeleteMovie(db *repository.Database) error {

	err := postgres.DeleteSingleRecord(db.Pdb.DB, `id = ?`, m, m.ID)

	if err != nil {
		return err
	}
	return nil
}

// custom type for date
func (d *Date) Scan(value interface{}) error {
	var t time.Time
	var err error
	switch value.(type) {
	case time.Time:
		t = value.(time.Time)
		t, err = time.Parse("2006-01-02", t.Format("2006-01-02"))
		if err != nil {
			return err
		}
	case []byte:
		t, err = time.Parse("2006-01-02", string(value.([]byte)))
		if err != nil {
			return err
		}
	default:
		return errors.New("failed to scan date")
	}
	*d = Date(t)
	return nil
}

func (d Date) Value() (driver.Value, error) {
	valueString, err := json.Marshal(d)
	return valueString, err
}

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var t string
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	dt, err := time.Parse("2006-01-02", t)
	if err != nil {
		return err
	}
	*d = Date(dt)
	return nil
}

// Custom type For Duration
func (d *Duration) Scan(value interface{}) error {
	switch value.(type) {
	case int64:
		*d = Duration(value.(int64))
	case []byte:
		intD, err := strconv.Atoi(string(value.([]byte)))
		if err != nil {
			return err
		}
		*d = Duration(uint8(intD))
	default:
		return errors.New("failed to scan duration")
	}
	return nil
}

func (d Duration) Value() (driver.Value, error) {
	valueString := uint8(d)
	return valueString, nil
}

// marshal json
func (d Duration) MarshalJSON() ([]byte, error) {
	hour := d / 60
	min := d % 60
	timeDuration := fmt.Sprintf("%dhr %dmin", hour, min)
	return json.Marshal(timeDuration)
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var t string
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	dt, err := strconv.Atoi(t)
	if err != nil {
		return err
	}
	*d = Duration(dt)
	return nil
}
