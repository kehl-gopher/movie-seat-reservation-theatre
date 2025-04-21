package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"gorm.io/gorm"
)

type Date time.Time // define custom date for release date

type Movie struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"type:text;not null;index"`
	Synopsis     string         `json:"synopsis" gorm:"type:text;not;index"`
	BackDropPath string         `json:"backdrop_path" gorm:"type:text"`
	PosterPath   string         `json:"poster_path" gorm:"type:text"`
	ReleaseDate  Date           `json:"release_date" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Genres       []Genre        `json:"genres" gorm:"many2many:movie_genres;"`
}

type Genre struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:text;not null;unique;index"`
}

func (d *Date) Scan(value interface{}) error {
	var t time.Time
	if err := json.Unmarshal(value.([]byte), &t); err != nil {
		fmt.Println(err.Error())
		return err
	}
	dt, err := time.Parse("2006-01-02", t.Format("2006-01-02"))

	if err != nil {
		return err
	}
	*d = Date(dt)
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
	var t time.Time
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	dt, err := time.Parse("2006-01-02", t.Format("2006-01-02"))
	if err != nil {
		return err
	}
	*d = Date(dt)
	return nil
}

func (m *Movie) BeforeCreate(tx *gorm.DB) error {
	m.ID = utility.GenerateUUID()
	return nil
}

func (m *Movie) CreateMovie(db *repository.Database) error {
	if err := postgres.Create(db.Pdb.DB, m); err != nil {
		return err
	}
	return nil
}

// TODO: get all movie relations
func (m *Movie) GetMovieByID(db *repository.Database, id string) (*Movie, error) {
	movie := &Movie{}
	return movie, nil
}

// TODO: get all movies with paginations paginations
func (m *Movie) GetAllMovies(db *repository.Database, page, limit int) ([]Movie, error) {
	var movies []Movie
	return movies, nil
}

// TODO: update movies
func (m *Movie) UpdateMovie(db *repository.Database, id string) (*Movie, error) {
	movie := &Movie{}
	return movie, nil
}

// TODO: delete movie
func (m *Movie) DeleteMovie(db *repository.Database, id string) error {
	return nil
}
