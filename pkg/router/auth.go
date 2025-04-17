package router

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/user"
)

func AuthRoutes(router *gin.Engine, config *env.Config, DB *repository.Database) {
	exp_in, _ := strconv.Atoi(config.EXPIRES_IN)
	userR := user.UserBase{DB: DB, ExpiresIn: int64(exp_in), SecretKey: []byte(config.SECRET_KEY)}
	userRoutes := router.Group(fmt.Sprintf("%s/%s", config.BASEURL, "auth"))
	{
		userRoutes.POST("/register", userR.UserSignUp)
		userRoutes.POST("/login", userR.UserSignIn)
	}
}
