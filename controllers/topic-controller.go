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
	if !ok {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Name is required"))
		return
	}

	topicInserted := &validations.CreateTopicWithUSerIDInput{
		Name:   input.Name,
		UserID: userId.(primitive.ObjectID),
	}
	topic, err := TopicModel.CreateTopic(topicInserted)
	if err != nil {
		ctx.Error(err)
		return
	}
	err = UserModel.UpdateTopicID(userId.(primitive.ObjectID), topic.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ecgController.CreateECGData(&topic.ID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"topicId": topic.ID})
	ctx.Abort()
}

func (t TopicController) GetTopic(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	retrievedData := &validations.GetTopicByIdInput{
		TopicID: topicId,
		UserID:  userId.(primitive.ObjectID),
	}
	topic, err := TopicModel.GetTopicById(retrievedData)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"topic": topic})
	ctx.Abort()
}

func (t TopicController) GetTopicForDoctor(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	topic, err := TopicModel.GetTopicByIdForDoctor(&topicId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"topic": topic})
	ctx.Abort()
}

func (t TopicController) UpdateTopicAnalyze(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	var updatedTopic validations.UpdateAnalyzeData
	if err := ctx.ShouldBindJSON(&updatedTopic); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	updatedData := validations.UpdateAnalyzeData{
		TopicID:  topicId,
		Analyzed: updatedTopic.Analyzed,
	}

	err = TopicModel.UpdateTopicAnalyzeData(&updatedData)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}

func (t TopicController) UpdateTopicAnalyzeComment(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	var updatedTopic validations.UpdateAnalyzeCommentInput
	userID, _ := ctx.Get("user_id")
	if err := ctx.ShouldBindJSON(&updatedTopic); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	user, err := UserModel.GetUserByID(userID.(primitive.ObjectID))
	if err != nil {
		ctx.Error(err)
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
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}

func (t TopicController) DeleteTopicAnalyzeComment(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("topicId"))
	userId, _ := ctx.Get("user_id")
	if err == nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid topic"))
		return
	}
	var deletedTopic validations.DeleteAnalyzeCommentInput
	if err := ctx.ShouldBindJSON(&deletedTopic); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	topicToDelete := validations.DeleteAnalyzeComment{
		TopicID:      topicId,
		CommentIndex: deletedTopic.CommentIndex,
		DoctorID:     userId.(primitive.ObjectID),
	}
	err = TopicModel.DeleteTopicAnalyzeComment(&topicToDelete)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}

func (t TopicController) DeleteTopic(ctx *gin.Context) {
	topicId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid topic ID"))
		return
	}
	err = TopicModel.DeleteTopic(topicId)
	if err != nil {
		ctx.Error(err)
		return
	}
	err = UserModel.DeleteTopicID(topicId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (t TopicController) UpdateTopicRecordTime(ctx *gin.Context) {
	var input validations.UpdateRecordTimeVal
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	newRecordTime := validations.UpdateRecordTimeVal{
		UserID:     userId.(primitive.ObjectID),
		TopicID:    input.TopicID,
		RecordTime: input.RecordTime,
	}
	err := TopicModel.UpdateRecordTime(&newRecordTime)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}

func (t TopicController) PredictionECGPlot(ctx *gin.Context) {
	var input validations.Prediction
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.Error(errorInstance.ReturnError(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	err := ecgController.ResampleECGData(input.TopicID)
	if err != nil {
		ctx.Error(err)
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
			ctx.Error(err)
			return
		}
		_, err = PredictionModel.CreatePrediction(prediction)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
	ctx.Abort()
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
			log.Println(err)
		}
		_, err = PredictionModel.CreatePrediction(prediction)
		if err != nil {
			log.Println(err)
		}
	}
}
