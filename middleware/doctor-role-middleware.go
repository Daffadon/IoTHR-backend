package middleware

import (
	"IoTHR-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DoctorMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token != "" {
		token = token[7:]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
			ctx.Abort()
			return
		}
		role := claims.Role
		if role == "doctor" {
			ctx.Next()
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		ctx.Abort()
		return
	}
}
