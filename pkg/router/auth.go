package router

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller/auth"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/middleware"
)

func AuthRoutes(router *gin.Engine, config *env.Config, DB *repository.Database) {
	exp_in, _ := strconv.Atoi(config.EXPIRES_IN)
	userR := auth.UserBase{DB: DB, ExpiresIn: int64(exp_in), SecretKey: []byte(config.SECRET_KEY)}
	userRoutes := router.Group(fmt.Sprintf("%s/%s", config.BASEURL, "auth"))
	{
		userRoutes.POST("/register", userR.SignUp)
		userRoutes.POST("/login", userR.SignIn)
	}
	adminP := router.Group(fmt.Sprintf("%s/auth", config.BASEURL))
	{
		adminP.POST("/admin/register", userR.AdminSignUp)
		adminP.POST("/admin/signin", userR.SignIn)
	}
	routesP := router.Group(fmt.Sprintf("%s", config.BASEURL), middleware.AuthMiddleWare(config.SECRET_KEY, DB))
	{
		routesP.POST("/logout", userR.SignOut)
	}

}
