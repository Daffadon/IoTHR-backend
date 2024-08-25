package server

import (
	"IoTHR-backend/controllers"
	"IoTHR-backend/middleware"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)	
	router := gin.New()
	health := new(controllers.HealthController)
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			if params.Method == "OPTIONS" {
				return ""
			}
			return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s |%-7s %#v\n%s",
				params.TimeStamp.Format(time.RFC1123),
				params.StatusCode,
				params.Latency,
				params.ClientIP,
				params.Method,
				params.Path,
				params.ErrorMessage,
			)
		},
		Output: gin.DefaultWriter,
	}))

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorMiddleware())
	router.Use(gin.Recovery())

	webSocket := new(controllers.WebsocketController)
	router.GET("/health", health.Status)
	router.GET("/ecg", webSocket.UpdateECGPlot)

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
			profile.Use(middleware.AuthMiddleware)
			user := new(controllers.UserController)
			profile.GET("", user.GetProfile)
			profile.GET("/history", user.GetHistory)
		}
		topic := v1.Group("/topic")
		{
			topic.Use(middleware.AuthMiddleware)
			topicController := new(controllers.TopicController)
			topic.POST("/create", topicController.CreateTopic)
			topic.GET("/:id", topicController.GetTopic)
			topic.POST("/prediction", topicController.PredictionECGPlot)
			topic.PATCH(("/record-time"), topicController.UpdateTopicRecordTime)
			topic.DELETE("/:id", topicController.DeleteTopic)
		}
		prediction := v1.Group("/prediction")
		{
			prediction.Use(middleware.AuthMiddleware)
			predictionController := new(controllers.PredictionController)
			prediction.GET(":id", predictionController.GetPredictionList)
			prediction.GET("id/:id", predictionController.GetPredictionById)
		}
		doctor := v1.Group("/doctor")
		{
			doctor.Use(middleware.AuthMiddleware)
			doctor.Use(middleware.DoctorMiddleware)
			user := new(controllers.UserController)
			topic := new(controllers.TopicController)
			doctor.GET("/user", user.GetUsers)
			doctor.GET("/user/:userId", user.GetUser)
			doctor.GET("/history/:userId", user.GetUserHistory)
			doctor.GET("/topic/:topicId", topic.GetTopicForDoctor)
			doctor.PATCH("/topic/analyze/:topicId", topic.UpdateTopicAnalyze)
			doctor.PATCH("/topic/analyze-comment/:topicId", topic.UpdateTopicAnalyzeComment)
			doctor.DELETE("/topic/analyze-comment/:topicId", topic.DeleteTopicAnalyzeComment)
		}
	}
	return router
}
