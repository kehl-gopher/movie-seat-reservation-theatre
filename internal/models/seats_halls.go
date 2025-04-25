package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Halls struct {
	ID    string  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name" gorm:"not null;unique"`
	Seats []Seats `json:"seats" gorm:"foreignKey:HallID;references:ID"`
}

type SeatStatus string

const (
	Available SeatStatus = "available"
	Held      SeatStatus = "held"
	Booked    SeatStatus = "booked"
)

func (s *SeatStatus) String() string {
	switch *s {
	case Available:
		return "available"
	case Held:
		return "held"
	case Booked:
		return "booked"
	default:
		return "unknown"
	}
}

type Seats struct {
	ID     string     `json:"id" gorm:"primaryKey;not null;autoIncrement"`
	Row    string     `json:"row" gorm:"not null;uniqueIndex:idx_row_number_hall_id"`
	Number int        `json:"number" gorm:"not null;uniqueIndex:idx_row_number_hall_id"`
	Status SeatStatus `json:"status" gorm:"not null;type:enum('available','held','booked');default:'available'"`
	HeldAt *time.Time `json:"held_at" gorm:"default:NULL"` // used for both optimistic locking and held status time minuete which cannot be > 10 minutes
	HallID string     `json:"hall_id" gorm:"not null;uniqueIndex:idx_row_number_hall_id"`
	Halls  Halls      `json:"halls" gorm:"foreignKey:HallID;references:ID"`
}

func (s SeatStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SeatStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	switch status {
	case string(Available), string(Held), string(Booked):
		*s = SeatStatus(status)
		return nil
	default:
		return fmt.Errorf("invalid seat status: %s", status)
	}
}

func (s *SeatStatus) Scan(value interface{}) error {
	if status, ok := value.(string); ok {
		*s = SeatStatus(status)
		return nil
	}
	return fmt.Errorf("invalid seat status: %v", value)
}
func (s SeatStatus) Value() (driver.Value, error) {
	return s, nil
}
