package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"gorm.io/gorm"
)

type Halls struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(20);not null;unique;index"`
	Seats     []Seats        `json:"seats" gorm:"foreignKey:HallID;references:ID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreatedAt"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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
	ID     string     `json:"id" gorm:"primaryKey;not null"`
	Row    string     `json:"row" gorm:"type:varchar(5);not null;uniqueIndex:idx_row_number_hall_id"`
	Number int        `json:"number" gorm:"not null;uniqueIndex:idx_row_number_hall_id"`
	Status SeatStatus `json:"status" gorm:"not null;type:seat_status;default:'available';index"`
	HeldAt *time.Time `json:"held_at" gorm:"default:NULL"`
	HallID string     `json:"hall_id" gorm:"not null;uniqueIndex:idx_row_number_hall_id;index"`
	Halls  Halls      `json:"halls" gorm:"foreignKey:HallID;references:ID"`
}

const MaxCinemaCapacity = 5000 // maximum cinema capacity

func generateRowsAndSeats(hall Halls, rowLabels, seatCount int) ([]Seats, error) {

	if rowLabels*seatCount > MaxCinemaCapacity {
		return nil, fmt.Errorf("theater exceeds maximum allowed capacity of %d seat", MaxCinemaCapacity)
	}
	seats := make([]Seats, rowLabels*seatCount)
	for r := 0; r < rowLabels; r++ {
		label := generateSeatLabel(r)
		for s := 1; s < seatCount+1; s++ {
			seat := Seats{
				ID:     utility.GenerateUUID(),
				Row:    label,
				Number: s,
				Status: Available,
				HallID: hall.ID,
				Halls:  hall,
			}
			seats[r*seatCount+s-1] = seat
		}
	}
	return seats, nil
}

func generateSeatLabel(ind int) string {
	var seatLabel strings.Builder
	chars := []byte{}
	for ind >= 0 {
		remainder := ind % 26
		chars = append(chars, byte(remainder+'A'))
		ind = ind/26 - 1
	}

	// rearrange characters...
	for i := len(chars) - 1; i >= 0; i-- {
		seatLabel.WriteByte(chars[i])
	}
	return seatLabel.String()
}

func (h *Halls) BeforeCreate(tx *gorm.DB) error {
	h.ID = utility.GenerateUUID()
	return nil
}

func (h *Halls) CreateHallSeat(db *repository.Database, config *env.Config, rowCount, seatCount int) (*Halls, int, error) {
	err := postgres.Create(db.Pdb.DB, h)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, http.StatusBadRequest, fmt.Errorf("hall name already exists %v", err.Error())
		}
		fmt.Println(err)
		return nil, http.StatusInternalServerError, err
	}

	seats, err := generateRowsAndSeats(*h, rowCount, seatCount)

	if err != nil {
		if err.Error() == fmt.Sprintf("theater exceeds maximum allowed capacity of %d seat", MaxCinemaCapacity) {
			return nil, http.StatusBadRequest, err
		}
		return nil, http.StatusInternalServerError, err
	}

	err = postgres.CreateMany(db.Pdb.DB, seats)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	hall := &Halls{}

	tx := postgres.Preload(db.Pdb.DB, hall, `Seats`)
	err = postgres.SelectById(tx, h.ID, &Halls{}, hall)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return hall, http.StatusCreated, err
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
