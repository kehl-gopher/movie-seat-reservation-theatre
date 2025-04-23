package movies

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/service/movies"
)

func (m *MovieBase) GetGenres(c *gin.Context) {

	genres, statusCode, err := movies.GetAllGenres(m.DB)

	if err != nil {
		fmt.Println(err)
		resp := utility.BuildErrorResponse(statusCode, err, "", http.StatusText(http.StatusInternalServerError))

		c.AbortWithStatusJSON(statusCode, resp)
		return
	}

	resp := utility.BuildSuccessResponse(statusCode, "", genres, nil)
	c.JSON(statusCode, resp)
}
