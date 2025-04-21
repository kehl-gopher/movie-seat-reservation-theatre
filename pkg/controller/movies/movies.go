package movies

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
)

type MovieBase struct {
	DB     *repository.Database
	Config *env.Config
}

func (m *MovieBase) CreateMovie(c *gin.Context) {
}
