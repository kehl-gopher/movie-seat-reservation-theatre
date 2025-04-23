package utility

import "github.com/gin-gonic/gin"

func GetParams(c *gin.Context, key string) string {
	return c.Param(key)
}