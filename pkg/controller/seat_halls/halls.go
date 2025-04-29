package seathalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	seathalls "github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/service/seat_halls"
)

type SeatHallBase struct {
	DB     *repository.Database
	Config *env.Config
}

func (h *SeatHallBase) CreateSeatHall(c *gin.Context) {

	var SeatHalls struct {
		HallName      string `json:"hall_name"`
		NumberOfRows  int    `json:"number_of_rows"`
		NumberOfSeats int    `json:"number_of_seats"`
	}
	err := c.ShouldBindJSON(&SeatHalls)
	if err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "bad request", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	data, statusCode, err := seathalls.CreateHallSeat(h.DB, h.Config, SeatHalls.HallName, SeatHalls.NumberOfRows, SeatHalls.NumberOfSeats)
	if err != nil {
		if err.Error() == "validation error" {
			v := err.(*utility.ValidationError)
			resp := utility.ValidationErrorResponse(v.Errors, v)
			c.JSON(http.StatusUnprocessableEntity, resp)
			return
		}
		resp := utility.BuildErrorResponse(statusCode, err, "", http.StatusText(statusCode))
		c.JSON(statusCode, resp)
		return
	}
	resp := utility.BuildSuccessResponse(statusCode, "created successfully", data, nil)
	c.JSON(statusCode, resp)
}

func (h *SeatHallBase) GetAllHalls(c *gin.Context) {
	data, statusCode, err := seathalls.GetAllHalls(h.DB, h.Config)
	if err != nil {
		resp := utility.BuildErrorResponse(statusCode, err, "", http.StatusText(statusCode))
		c.JSON(statusCode, resp)
		return
	}
	resp := utility.BuildSuccessResponse(statusCode, "created successfully", data, nil)
	c.JSON(statusCode, resp)
}

// TODO: Implement admin updating seat hall later stupid fuck
func (h *SeatHallBase) UpdateSeatHall(c *gin.Context) {

	hallID := utility.GetParams(c, "hallID")

	if _, err := uuid.Parse(hallID); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid uuid ID parsed", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var SeatHall struct {
		HallName      *string `json:"hall_name"`
		RowNumber     *string `json:"row_number"`
		NumberOfSeats *int    `json:"number_of_seats"`
	}

	err := c.ShouldBindJSON(&SeatHall)
	if err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "bad error response", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

}
