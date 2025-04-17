package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

func AuthMiddleWare(secret_key string, db *repository.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorisation := ctx.Request.Header.Get("Authorization")
		auths := strings.Split(authorisation, " ")
		if len(auths) != 2 {
			resp := utility.UnAuthorizedResponse(http.StatusBadRequest, "invalid token", http.StatusText(http.StatusBadRequest))
			ctx.JSON(http.StatusBadRequest, resp)
			ctx.Abort()
			return
		}
		if auths[0] != "Bearer" {
			resp := utility.UnAuthorizedResponse(http.StatusBadRequest, "invalid token", http.StatusText(http.StatusBadRequest))
			ctx.JSON(http.StatusBadRequest, resp)
			ctx.Abort()
			return
		}

		userClaim := utility.AccessTokenClaim{SecretKey: []byte(secret_key)}
		userId, roleID, err := userClaim.ExtractClaims(auths[1])
		if err != nil {
			resp := utility.UnAuthorizedResponse(http.StatusBadRequest, err.Error(), http.StatusText(http.StatusBadRequest))
			ctx.JSON(http.StatusBadRequest, resp)
			ctx.Abort()
			return
		}

		// set user object in context
		ctx.Set("roleID", models.RoleIDs(roleID))
		user := models.Users{ID: userId}
		u, err := user.PreloadUserRole(db, models.RoleIDs(roleID))

		if err != nil {
			resp := utility.UnAuthorizedResponse(http.StatusBadRequest, err.Error(), http.StatusText(http.StatusBadRequest))
			ctx.JSON(http.StatusBadRequest, resp)
			ctx.Abort()
			return
		}
		ctx.Set("user", *u)
	}
}
