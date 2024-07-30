package middleware

import (
	"IoTHR-backend/models"
	"IoTHR-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UserModels = new(models.User)

func AuthMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token != "" {
		token = token[7:]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
			ctx.Abort()
			return
		}
		loggedOut := UserModels.IsLoggedOut(token)
		if !loggedOut {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserId)
		ctx.Set("role", claims.Role)
		ctx.Next()
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
	ctx.Abort()
}
