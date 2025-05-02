package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
)

type BookingBase struct {
	DB     *repository.Database
	Config *env.Config
}

func (b *BookingBase) CreateBooking(c *gin.Context) {

}
