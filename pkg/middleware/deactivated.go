package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
)

// Handle deactivated user routes
func DeactivatedMiddleware(db *repository.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
	}
}
