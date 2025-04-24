package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		return
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

func (m *MovieBase) GetMovie(c *gin.Context) {
	movie_id := utility.GetParams(c, "movieId")
	if _, err := uuid.Parse(movie_id); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid movie id", http.StatusText(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	movie, statusCode, err := movies.GetMovieByID(m.DB, movie_id)

	if err != nil {
		resp := utility.BuildErrorResponse(statusCode, err, "error getting movie", http.StatusText(statusCode))
		c.AbortWithStatusJSON(statusCode, resp)
		return
	}

	resp := utility.BuildSuccessResponse(statusCode, "movie found", movie, nil)
	c.JSON(statusCode, resp)

}

func (m *MovieBase) GetMovies(c *gin.Context) {

	limit, offset := utility.GetPaginationParams(c)

	movies, statusCode, pag, err := movies.GetAllMovies(m.DB, limit, offset, m.Config)

	if err != nil {
		resp := utility.BuildErrorResponse(statusCode, err, "error getting movies", http.StatusText(statusCode))
		c.AbortWithStatusJSON(statusCode, resp)
		return
	}

	resp := utility.BuildSuccessResponse(statusCode, "movies found", movies, pag)

	c.JSON(statusCode, resp)

}

func (m *MovieBase) UpdateMovie(c *gin.Context) {
	movie_id := utility.GetParams(c, "movieId")
	if _, err := uuid.Parse(movie_id); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid movie id", http.StatusText(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	mov := &movies.UpdateMovieReq{}

	if err := c.ShouldBindJSON(mov); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "bad request sent", http.StatusText(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	statusCode, err := movies.UpdateMovie(m.DB, movie_id, mov, m.Config)

	if err != nil {
		if err.Error() == "validation error" {
			v := err.(*utility.ValidationError)
			resp := utility.ValidationErrorResponse(v.Errors, v)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, resp)
			return
		}
		resp := utility.BuildErrorResponse(statusCode, err, "error updating movie", http.StatusText(statusCode))
		c.AbortWithStatusJSON(statusCode, resp)
		return
	}
	resp := utility.BuildSuccessResponse(statusCode, "movie updated successfully", nil, nil)
	c.JSON(http.StatusOK, resp)
}

func (m *MovieBase) DeleteMovie(c *gin.Context) {
	movie_id := utility.GetParams(c, "movieId")
	if _, err := uuid.Parse(movie_id); err != nil {
		resp := utility.BuildErrorResponse(http.StatusBadRequest, err, "invalid movie id", http.StatusText(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	statusCode, err := movies.DeleteMovie(m.DB, movie_id)

	if err != nil {
		resp := utility.BuildErrorResponse(statusCode, err, "error deleting movie", http.StatusText(statusCode))
		c.AbortWithStatusJSON(statusCode, resp)
		return
	}
	resp := utility.BuildSuccessResponse(statusCode, "movie deleted successfully", nil, nil)
	c.JSON(http.StatusOK, resp)
}
