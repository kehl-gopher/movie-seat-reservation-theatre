package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is an admin
		roleID, exists := c.Get("roleID")

		id := roleID.(models.RoleIDs)
		if !exists || id.String() != models.Admin {
			resp := utility.UnAuthorizedResponse("user not an admin", http.StatusText(http.StatusUnauthorized))
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		// Continue to the next handler
		c.Next()
	}
}
