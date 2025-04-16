package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/user"
)

func AuthRoutes(router *gin.Engine, config *env.Config, DB *repository.Database) {

	userR := user.UserBase{DB: DB}
	userRoutes := router.Group(fmt.Sprintf("%s/%s", config.BASEURL, "auth"))
	{
		fmt.Println(userRoutes.BasePath())
		userRoutes.POST("/register", userR.UserSignUp)
	}

}
