package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	seathalls "github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/seat_halls"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/middleware"
)

func SeatHallRoutes(r *gin.Engine, db *repository.Database, config *env.Config) {
	seatHall := r.Group(fmt.Sprintf("%s/admin/seat-hall", config.BASEURL), middleware.AuthMiddleWare(config.SECRET_KEY, db), middleware.AuthAdmin())
	hallBase := seathalls.SeatHallBase{DB: db, Config: config}
	{
		seatHall.POST("/", hallBase.CreateSeatHall)
		seatHall.GET("/", hallBase.GetAllHalls)
		seatHall.GET("/:hallId", hallBase.GetHall)
		seatHall.DELETE("/:hallId", hallBase.DeleteHall)
	}
}
