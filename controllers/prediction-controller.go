package controllers

import (
	"IoTHR-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PredictionController struct{}

var PredictionModel = new(models.Prediction)

func (p PredictionController) GetPredictionList(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	predictions, err := PredictionModel.GetPredictionByTopicId(&topicId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(200, gin.H{"data": predictions})
	ctx.Abort()
}

func (p PredictionController) GetPredictionById(ctx *gin.Context) {
	predictionId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	prediction, err := PredictionModel.GetPredictionByPredictionId(&predictionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(200, gin.H{"data": prediction})
	ctx.Abort()
}
