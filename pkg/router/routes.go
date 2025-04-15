package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/router/app"
)

type Application struct {
}

func Router(config *env.Config, DB *repository.Database) error {
	router := gin.Default()

	// application status
	app.RunApp(router, config)

	router.Run() // listen and serve on port 8000
	return nil
}
