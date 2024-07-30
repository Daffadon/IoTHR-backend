package controllers

import (
	"IoTHR-backend/models"
	"IoTHR-backend/validations"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicController struct{}

var TopicModel = new(models.Topic)

func (t TopicController) CreateTopic(ctx *gin.Context) {
	var input validations.CreateTopicInput
	userId, ok := ctx.Get("user_id")
	if ok {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
		}
		topicInserted := &validations.CreateTopicWithUSerIDInput{
			Name:   input.Name,
			UserID: userId.(primitive.ObjectID),
		}
		topic, err := TopicModel.CreateTopic(topicInserted)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		err = UserModel.UpdateTopicID(userId.(primitive.ObjectID), topic.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}		
		ctx.JSON(http.StatusOK, gin.H{"topicId": topic.ID})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func (t TopicController) GetTopic(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err == nil {
		retrievedData := &validations.GetTopicByIdInput{
			TopicID: topicId,
			UserID:  userId.(primitive.ObjectID),
		}
		topic, err := TopicModel.GetTopicById(retrievedData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"topic": topic})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	ctx.Abort()
}

func (t TopicController) UpdateECGPlotTopic(ctx *gin.Context) {
	var input validations.InsertECGDataInput
	userId, ok := ctx.Get("user_id")
	if ok {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
		}
		dataToUpdate := validations.UpdateECGInput{
			TopicID: input.TopicID,
			UserID:  userId.(primitive.ObjectID),
			ECGPlot: input.ECGPlot,
		}
		err := TopicModel.UpdateECGdata(&dataToUpdate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func (t TopicController) PredictionECGPlot(ctx *gin.Context) {
	var input validations.Prediction
	userId, ok := ctx.Get("user_id")
	if ok {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
		}
		newRecordTime := &validations.UpdateRecordTimeVal{
			UserID:     userId.(primitive.ObjectID),
			TopicID:    input.TopicID,
			RecordTime: input.RecordTime,
		}
		err := TopicModel.UpdateRecordTime(newRecordTime)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		feature := []string{"RR-Interval", "Morphology", "Wavelet"}
		for _, f := range feature {
			dataToPredict := &validations.ECGPredictionInput{
				TopicID: input.TopicID,
				UserID:  userId.(primitive.ObjectID),
				Feature: f,
			}
			prediction, err := TopicModel.ECGPrediction(dataToPredict)
			if err != nil {
				log.Fatal(err)
			}
			_, err = PredictionModel.CreatePrediction(prediction)
			if err != nil {
				log.Fatal(err)
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
