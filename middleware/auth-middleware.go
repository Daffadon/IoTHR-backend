package middleware

import (
	"IoTHR-backend/errors"
	"IoTHR-backend/models"
	"IoTHR-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UserModels = new(models.User)
var errorInstance = new(errors.ErrorInstance)

func AuthMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	token = token[7:]
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	loggedOut := UserModels.IsLoggedOut(token)
	if !loggedOut {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	ctx.Set("user_id", claims.UserId)
	ctx.Set("role", claims.Role)
	ctx.Next()
}
