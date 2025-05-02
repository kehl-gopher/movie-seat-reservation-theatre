package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ShowTime time.Time

type Shows struct {
	ID        string          `json:"id" gorm:"primaryKey"`
	MovieID   string          `json:"movie_id" gorm:"not null"`
	HallID    string          `json:"hall_id" gorm:"not null"`
	StartDate Date            `json:"start_date" gorm:"type:Date;not null"`
	StartTime ShowTime        `json:"start_time" gorm:"type:Time WITHOUT TIME ZONE;not null"`
	EndTime   ShowTime        `json:"end_time" gorm:"type:Time WITHOUT TIME ZONE;not null"`
	Price     decimal.Decimal `json:"price" gorm:"type:Decimal(10, 2);not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	Movie     *Movie          `json:"movie,omitempty" gorm:"foreignKey:MovieID;references:ID"`
	Hall      *Halls          `json:"hall,omitempty" gorm:"foreignKey:HallID;references:ID"`
}

func (s *Shows) BeforeCreate(tx *gorm.DB) error {
	s.ID = utility.GenerateUUID()

	return nil
}
func (s *Shows) CreateMovieShows(db *repository.Database) error {

	h := Halls{ID: s.HallID}
	_, err := h.GetHall(db)
	if err != nil {
		return err
	}

	m := Movie{ID: s.MovieID}

	_, err = m.GetMovieByID(db, m.ID)
	if err != nil {
		return err
	}

	err = postgres.Create(db.Pdb.DB, s)

	if err != nil {
		return err
	}

	return nil

}

func (t *ShowTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*t = ShowTime(v)
		return nil
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return fmt.Errorf("failed to parse ShowTime from string: %v", err)
		}
		*t = ShowTime(parsed)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into ShowTime", value)
	}
}

func (s *Shows) GetAllShows(db *repository.Database) ([]Shows, error) {
	shows := []Shows{}
	tx := postgres.Preload(db.Pdb.DB, s, `Movie`, `Hall`)
	err := postgres.SelectAllRecords(tx, "desc", "created_at", s, &shows)

	if err != nil {
		return nil, err
	}
	return shows, nil
}

func (s ShowTime) Value() (driver.Value, error) {
	return time.Time(s).Format("15:04"), nil
}

func (s ShowTime) MarshalJSON() ([]byte, error) {
	t := time.Time(s)
	return json.Marshal(t.Format("15:04"))
}

func (s *ShowTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	tie := string(data)
	tie = strings.Trim(tie, "\"")
	t, err := time.Parse("15:04", tie)
	if err != nil {
		return err
	}
	*s = ShowTime(t)
	return nil
}
