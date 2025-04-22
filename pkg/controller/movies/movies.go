package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/service/movies"
)

type MovieBase struct {
	DB     *repository.Database
	Config *env.Config
}

func (m *MovieBase) CreateMovie(c *gin.Context) {

	mov := &movies.MovieReq{}

	if err := c.ShouldBindJSON(mov); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "bad request sent", http.StatusText(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
	}

	statusCode, err := movies.CreateMovies(m.DB, mov, m.Config)

	if err != nil {
		if err.Error() == "validation error" {
			v := err.(*utility.ValidationError)
			resp := utility.ValidationErrorResponse(v.Errors, v)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, resp)
			return
		}
		resp := utility.BuildErrorResponse(statusCode, err, "error creating movie", http.StatusText(statusCode))
		c.AbortWithStatusJSON(statusCode, resp)
		return
	}
	resp := utility.BuildSuccessResponse(http.StatusCreated, "movie created successfully", nil, http.StatusText(http.StatusCreated))
	c.JSON(http.StatusCreated, resp)
}
