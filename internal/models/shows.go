package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ShowTime time.Time

type Shows struct {
	ID        string          `json:"id" gorm:"primaryKey"`
	MovieID   string          `json:"movie_id" gorm:"not null"`
	HallID    string          `json:"hall_id" gorm:"not null"`
	StartDate Date            `json:"start_date" gorm:"type:Date;not null"`
	StartTime ShowTime        `json:"start_time" gorm:"type:Time;not null"`
	EndTime   ShowTime        `json:"end_time" gorm:"type:Time;not null"`
	Price     decimal.Decimal `json:"price" gorm:"not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	Movie     Movie           `json:"movie" gorm:"foreignKey:MovieID;references:ID"`
	Hall      Halls           `json:"hall" gorm:"foreignKey:HallID;references:ID"`
}

func (s *ShowTime) Scan(value interface{}) error {

	if t, ok := value.(time.Time); ok {
		t, err := time.Parse("15:04:05", t.Format("15:04:05"))
		if err != nil {
			return err
		}
		*s = ShowTime(t)
		return nil
	} else if t, ok := value.([]byte); ok {
		t, err := time.Parse("15:04:05", string(t))
		if err != nil {
			return err
		}
		*s = ShowTime(t)
		return nil
	}
	return fmt.Errorf("failed to scan ShowTime: %v", value)
}

