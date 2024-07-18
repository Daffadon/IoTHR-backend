package controllers

import (
	"IoTHR-backend/validations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (u UserController) GetProfile(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if ok {
		if user, err := UserModel.GetUserByID(userID.(uint)); err == nil {
			userData := &validations.Profile{
				Email: user.Email, Fullname: user.Fullname,
			}
			ctx.JSON(http.StatusOK, gin.H{"data": userData})
			ctx.Abort()
			return
		}
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
}
