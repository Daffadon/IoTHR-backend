package middleware

import (
	errInstance "IoTHR-backend/errors"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
			if err != nil {
				switch e := err.(type) {
				case *errInstance.ErrorInstance:
					ctx.JSON(e.Code, gin.H{"error": e.Message})
					ctx.Abort()
				default:
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					ctx.Abort()
				}
			}
		}
	}
}
