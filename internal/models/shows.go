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
	Movie     Movie           `json:"movie" gorm:"foreignKey:MovieID;references:ID"`
	Hall      Halls           `json:"hall" gorm:"foreignKey:HallID;references:ID"`
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
func (s *ShowTime) Scan(value interface{}) error {

	if t, ok := value.(time.Time); ok {
		t, err := time.Parse("15:04", t.Format("15:04"))
		if err != nil {
			return err
		}
		*s = ShowTime(t)
		return nil
	} else if t, ok := value.([]byte); ok {
		t, err := time.Parse("15:04", string(t))
		if err != nil {
			return err
		}
		*s = ShowTime(t)
		return nil
	}
	return fmt.Errorf("failed to scan ShowTime: %v", value)
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
