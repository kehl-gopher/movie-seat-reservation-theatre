package shows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/service/shows"
	"github.com/shopspring/decimal"
)

type ShowBase struct {
	DB     *repository.Database
	Config *env.Config
}

func (s *ShowBase) CreateShows(c *gin.Context) {
	var Shows struct {
		MovieID   string          `json:"movie_id"`
		HallID    string          `json:"hall_id"`
		Price     decimal.Decimal `json:"price"`
		StartDate models.Date     `json:"start_date"`
		StartTime models.ShowTime `json:"start_time"`
		EndTime   models.ShowTime `json:"end_time"`
	}

	err := c.ShouldBindJSON(&Shows)
	if err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "bad error response", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	movieId, err := uuid.Parse(Shows.MovieID)
	if err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid uuid parse", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	hallId, err := uuid.Parse(Shows.HallID)
	if err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid uuid parse", http.StatusText(http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	statusCode, err := shows.CreateShow(s.DB, movieId.String(), hallId.String(), Shows.Price, Shows.StartTime, Shows.EndTime, Shows.StartDate)

	if err != nil {
		if err.Error() == "validation error" {
			v := err.(*utility.ValidationError)
			resp := utility.ValidationErrorResponse(v.Errors, v)
			c.JSON(statusCode, resp)
			return

		}
		resp := utility.BuildErrorResponse(statusCode, err, "", http.StatusText(statusCode))
		c.JSON(statusCode, resp)
		return
	}

	resp := utility.BuildSuccessResponse(statusCode, "show time created successfully", nil, nil)
	c.JSON(statusCode, resp)
}
