package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
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
			http.StatusBadRequest, err, "error",
			http.StatusText(http.StatusUnprocessableEntity),
		)
		ctx.JSON(http.StatusUnprocessableEntity, resp)
		return
	}
	roleId := 2 // role id for user
	statusCode, succResp, err := auth.UserRequestService(user, u.DB, models.RoleIDs(roleId), u.ExpiresIn, u.SecretKey)

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

	resp := utility.BuildSuccessResponse(statusCode, "success", succResp, nil)
	ctx.JSON(statusCode, resp)
}

func (u *UserBase) UserSignIn(ctx *gin.Context) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		resp := utility.BuildErrorResponse(
			http.StatusBadRequest, err, "error",
			http.StatusText(http.StatusUnprocessableEntity),
		)
		ctx.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	roleID := 2 // user role id
	statusCode, data, err := auth.AuthenticateUser(u.DB, user.Email, user.Password, u.SecretKey, u.ExpiresIn, models.RoleIDs(roleID))

	if err != nil {
		eresp := utility.BuildErrorResponse(statusCode, err, "error", http.StatusText(statusCode))
		ctx.JSON(statusCode, eresp)
		return
	}
	sresp := utility.BuildSuccessResponse(statusCode, "success", data, nil)
	ctx.JSON(statusCode, sresp)
}
