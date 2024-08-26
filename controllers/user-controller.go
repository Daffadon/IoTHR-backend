package controllers

import (
	"IoTHR-backend/validations"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct{}

func (u UserController) GetUsers(ctx *gin.Context) {
	users, err := UserModel.GetUsers()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func (u UserController) GetUser(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(ctx.Param("userId"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	user, err := UserModel.GetUserByID(userId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
	ctx.Abort()
}

func (u UserController) GetUserHistory(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(ctx.Param("userId"))
	if err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Invalid user ID"))
	}
	if user, err := UserModel.GetUserByID(userId); err == nil {
		topicList, err := TopicModel.GetUserTopicList(user.TopicID)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": topicList})
		ctx.Abort()
		return
	}
	ctx.Error(err)
}

func (u UserController) GetProfile(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid user"))
	}
	if user, err := UserModel.GetUserByID(userID.(primitive.ObjectID)); err == nil {
		userData := &validations.Profile{
			Email:     user.Email,
			Fullname:  user.Fullname,
			BirthDate: user.BirthDate,
		}
		ctx.JSON(http.StatusOK, gin.H{"data": userData})
		ctx.Abort()
		return
	}
	ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid user"))
}

func (u UserController) GetHistory(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid user"))
		return
	}
	if user, err := UserModel.GetUserByID(userID.(primitive.ObjectID)); err == nil {
		topicNames, err := TopicModel.GetTopicList(&validations.GetHistoryInput{TopicList: user.TopicID})
		if err != nil {
			ctx.Error(err)
			return
		}
		if (len(topicNames)) == 0 {
			ctx.Error(errorInstance.ReturnError(http.StatusNotFound, ("History not found")))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": topicNames})
		ctx.Abort()
		return
	}
}
