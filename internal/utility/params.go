package utility

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParams(c *gin.Context, key string) string {
	return c.Param(key)
}

func GetQuery(c *gin.Context, key string) string {
	return c.Query(key)
}

func GetPaginationParams(c *gin.Context) (uint, uint) {
	var pageInt, limitInt int
	page := c.Query("page")
	limit := c.Query("limit")

	if page == ""  {
		pageInt = 0
	}
	if limit == "" {
		limitInt = 10
	}
	pageInt, _ = strconv.Atoi(page)
	limitInt, _ = strconv.Atoi(limit)

	return uint(pageInt), uint(limitInt)
}
