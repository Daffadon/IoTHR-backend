package server

import (
	"IoTHR-backend/controllers"
	"IoTHR-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	health := new(controllers.HealthController)
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello World!"})
	})
	router.GET("/health", health.Status)
	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			auth := new(controllers.AuthController)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/logout", middleware.AuthMiddleware, auth.Logout)
		}
		profile := v1.Group("/profile")
		{
			user := new(controllers.UserController)
			profile.GET("", middleware.AuthMiddleware, user.GetProfile)
			profile.GET("/history", middleware.AuthMiddleware, user.GetHistory)
		}
		topic := v1.Group("/topic")
		{
			topic.Use(middleware.AuthMiddleware)
			topicController := new(controllers.TopicController)
			topic.POST("/create", topicController.CreateTopic)
			topic.PATCH("/ecg", topicController.UpdateECGPlotTopic)
			topic.GET("/:id", topicController.GetTopic)
			topic.PATCH("/prediction", topicController.PredictionECGPlot)
		}
		prediction := v1.Group("/prediction")
		{
			prediction.Use(middleware.AuthMiddleware)
			predictionController := new(controllers.PredictionController)
			prediction.GET(":id", predictionController.GetPredictionList)
			prediction.GET("id/:id", predictionController.GetPredictionById)
		}
	}
	return router
}
