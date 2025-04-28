package seathalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
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
}
