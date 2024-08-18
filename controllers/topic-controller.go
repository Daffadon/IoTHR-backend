package controllers

import (
	"IoTHR-backend/validations"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicController struct{}

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

		err = ecgController.CreateECGData(&topic.ID)
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

func (t TopicController) GetTopicForDoctor(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err == nil {
		topic, err := TopicModel.GetTopicByIdForDoctor(&topicId)
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

func (t TopicController) UpdateTopicAnalyze(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err == nil {
		var updatedTopic validations.UpdateAnalyzeData
		if err := ctx.ShouldBindJSON(&updatedTopic); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		updatedData := validations.UpdateAnalyzeData{
			TopicID:  topicId,
			Analyzed: updatedTopic.Analyzed,
		}

		err := TopicModel.UpdateTopicAnalyzeData(&updatedData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic"})
}

func (t TopicController) UpdateTopicAnalyzeComment(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err == nil {
		var updatedTopic validations.UpdateAnalyzeCommentInput
		if err := ctx.ShouldBindJSON(&updatedTopic); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		userID, _ := ctx.Get("user_id")
		user, err := UserModel.GetUserByID(userID.(primitive.ObjectID))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		updatedData := validations.UpdateAnalyzeComment{
			TopicID:    topicId,
			DoctorID:   userID.(primitive.ObjectID),
			DoctorName: user.Fullname,
			Comment:    []string{updatedTopic.Comment},
		}

		err = TopicModel.UpdateTopicAnalyzeComment(&updatedData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic"})
}

func (t TopicController) DeleteTopicAnalyzeComment(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	userId, _ := ctx.Get("user_id")
	if err == nil {
		var deletedTopic validations.DeleteAnalyzeCommentInput
		if err := ctx.ShouldBindJSON(&deletedTopic); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		topicToDelete := validations.DeleteAnalyzeComment{
			TopicID:      topicId,
			CommentIndex: deletedTopic.CommentIndex,
			DoctorID:     userId.(primitive.ObjectID),
		}
		err := TopicModel.DeleteTopicAnalyzeComment(&topicToDelete)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic"})
}

func (t TopicController) UpdateTopicRecordTime(ctx *gin.Context) {
	var input validations.UpdateRecordTimeVal
	userId, ok := ctx.Get("user_id")
	if ok {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
		}
		newRecordTime := validations.UpdateRecordTimeVal{
			UserID:     userId.(primitive.ObjectID),
			TopicID:    input.TopicID,
			RecordTime: input.RecordTime,
		}
		err := TopicModel.UpdateRecordTime(&newRecordTime)
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
		err := ecgController.ResampleECGData(input.TopicID)

		if err != nil {
			log.Fatalf("Error resampling ECG data: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
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

func (t TopicController) PredictionForWS(topicId, userId primitive.ObjectID) {
	feature := []string{"RR-Interval", "Morphology", "Wavelet"}
	for _, f := range feature {
		dataToPredict := &validations.ECGPredictionInput{
			TopicID: topicId,
			UserID:  userId,
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
}
