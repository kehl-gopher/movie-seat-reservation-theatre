package app

import "github.com/gin-gonic/gin"

type AppStat struct {
	AppName string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

// will define the application status but for now... it return pong
func (a *AppStat) AppVersioning(c *gin.Context) {
	c.JSON(200, gin.H{
		"app_name": a.AppName,
		"version":  a.Version,
		"env":      a.Env,
	})
}
