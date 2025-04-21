package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
)

func Router(config *env.Config, DB *repository.Database) error {
	router := gin.Default()

	// application status
	RunApp(router, config)
	AuthRoutes(router, config, DB)
	MovieRoutes(router, config, DB)

	lstn := fmt.Sprintf("%s:%s", config.APP_HOST, config.APP_PORT)
	err := router.Run(lstn)

	if err != nil {
		fmt.Println(err)
	}
	return err
}
