package shows

import (
	"net/http"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/shopspring/decimal"
)

func ValidateCreateShowEntry(startTime, endTime models.ShowTime, startDate models.Date, price decimal.Decimal) error {
	v := utility.NewValidationError()
	sd := time.Time(startDate).UTC()
	st := time.Time(startTime).UTC()
	now := time.Now().UTC().Truncate(time.Second)

	stDate := time.Date(sd.Year(), sd.Month(), sd.Day(), st.Hour(), st.Minute(), 0, 0, sd.Location())
	if stDate.Before(now) {
		v.AddValidationError("start_time", "start time must be in the past")
	}

	et := time.Time(endTime).UTC()
	endDate := time.Date(sd.Year(), sd.Month(), sd.Day(), et.Hour(), et.Minute(), 0, 0, sd.Location())

	if endDate.Before(now) || et.Before(st) {
		v.AddValidationError("end_time", "end time must be in the future or must be after start time")
	}
	if sd.Before(now) {
		v.AddValidationError("start_date", "start date should be present or be in the future")
	}
	dec0 := decimal.NewFromInt(0)
	if price.LessThanOrEqual(dec0) {
		v.AddValidationError("price", "price field should be greater than 0")
	}

	if v.CheckError() {
		return v
	}
	return nil
}

func CreateShow(movieId, hallId string, price decimal.Decimal, startTime, endTime models.ShowTime, startDate models.Date) (int, error) {

	err := ValidateCreateShowEntry(startTime, endTime, startDate, price)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}
	return 0, nil
}
