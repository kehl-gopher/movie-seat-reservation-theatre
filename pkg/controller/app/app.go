package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

type AppStat struct {
	AppName string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

func (a *AppStat) AppVersioning(c *gin.Context) {
	data := map[string]interface{}{
		"app_name": a.AppName,
		"version":  a.Version,
		"env":      a.Env,
	}
	resp := utility.BuildSuccessResponse(http.StatusOK, "pong", data, nil)
	c.JSON(http.StatusOK, resp)
}
