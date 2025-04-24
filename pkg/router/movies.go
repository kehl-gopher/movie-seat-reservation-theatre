package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/movies"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/middleware"
)

func MovieRoutes(router *gin.Engine, config *env.Config, DB *repository.Database) {
	movieBase := movies.MovieBase{DB: DB, Config: config}
	secret_key := []byte(config.SECRET_KEY)
	movieUrl := router.Group(fmt.Sprintf("%s", config.BASEURL))
	{
		movieUrl.POST("/movie", middleware.AuthMiddleWare(string(secret_key), DB), middleware.AuthAdmin(), movieBase.CreateMovie)
		movieUrl.GET("/movies", movieBase.GetMovies)
		movieUrl.GET("/movie/:movieId", movieBase.GetMovie)
		movieUrl.PATCH("/movie/:movieId", middleware.AuthMiddleWare(string(secret_key), DB), middleware.AuthAdmin(), movieBase.UpdateMovie)
		movieUrl.DELETE("/movie/:movieId", middleware.AuthMiddleWare(string(secret_key), DB), middleware.AuthAdmin(), movieBase.DeleteMovie)
	}
	// get movie genres
	genres := router.Group(fmt.Sprintf("%s/genres", config.BASEURL))
	{
		genres.GET("", movieBase.GetGenres)
	}

}
