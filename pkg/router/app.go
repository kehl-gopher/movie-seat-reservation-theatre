package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/controller"
)

func RunApp(router *gin.Engine, config *env.Config) {
	a := controller.AppStat{AppName: config.APP_NAME, Version: config.VERSION, Env: config.ENV}
	baseUrl := fmt.Sprintf("%s/ping", config.BASEURL)
	router.GET(baseUrl, a.AppVersioning)
}
