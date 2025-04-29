package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/shows"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/middleware"
)

func ShowRoutes(r *gin.Engine, db *repository.Database, config *env.Config) {

	showBase := shows.ShowBase{DB: db, Config: config}
	showRoutes := r.Group(fmt.Sprintf("%s/admin/show", config.BASEURL), middleware.AuthMiddleWare(config.SECRET_KEY, db), middleware.AuthAdmin())
	{
		showRoutes.POST("/", showBase.CreateShows)
	}
}
