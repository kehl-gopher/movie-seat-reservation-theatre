package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/service/auth"
)

type UserBase struct {
	DB        *repository.Database
	ExpiresIn int64
	SecretKey []byte
}

func (u *UserBase) UserSignUp(ctx *gin.Context) {
	var user auth.UserAuthRequest
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		resp := utility.BuildErrorResponse(
			http.StatusUnprocessableEntity,
			err,
			"validation error",
			http.StatusText(http.StatusUnprocessableEntity),
		)
		ctx.JSON(http.StatusUnprocessableEntity, resp)
	}
	statusCode, succResp, err := auth.UserRequestService(user, u.DB, 2, u.ExpiresIn, u.SecretKey)

	if err != nil {
		switch statusCode {
		case http.StatusUnprocessableEntity:
			v := err.(*utility.ValidationError)
			resp := utility.ValidationErrorResponse(v.Errors, v) // handle validation error logic
			ctx.JSON(statusCode, resp)
		default:
			resp := utility.BuildErrorResponse(statusCode, err, "error", http.StatusText(statusCode))
			ctx.JSON(statusCode, resp)
		}
		return
	}

	ctx.JSON(statusCode, succResp)
}
