package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PredictionController struct{}

func (p PredictionController) GetPredictionList(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	predictions, err := PredictionModel.GetPredictionByTopicId(&topicId)
	if err != nil {
		ctx.Error(err)
		return
	}
	if len(*predictions) == 0 {
		ctx.Error(errorInstance.ReturnError(http.StatusNotFound, "Prediction not found"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": predictions})
	ctx.Abort()
}

func (p PredictionController) GetPredictionById(ctx *gin.Context) {
	predictionId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid prediction ID"))
		return
	}
	prediction, err := PredictionModel.GetPredictionByPredictionId(&predictionId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": prediction})
	ctx.Abort()
}
