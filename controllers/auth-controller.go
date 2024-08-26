package controllers

import (
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (auth AuthController) Login(ctx *gin.Context) {
	var input validations.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid email or password"))
		return
	}
	user, err := UserModel.GetUser(&input)
	if err != nil {
		ctx.Error(err)
		return
	}
	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errCompare != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid email or password"))
		return
	}
	token, errJWT := utils.GenerateToken(user.ID, user.Role)
	if errJWT != nil {
		ctx.Error(errJWT)
		return
	}
	errUpdate := UserModel.UpdateTokenToUser(user.ID, token)
	if errUpdate != nil {
		ctx.Error(errUpdate)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": token})
	ctx.Abort()
}

func (auth AuthController) Register(ctx *gin.Context) {
	var input validations.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "invalid email or password"))
		return
	}

	if input.Password != input.ConfirmPassword {
		ctx.Error(errorInstance.ReturnError(http.StatusBadRequest, "Password and Confirm Password not match!"))
		return
	}

	if hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), 14); errHash != nil {
		ctx.Error(errorInstance.ReturnError(http.StatusInternalServerError, "Failed to hash password!"))
		return
	} else {
		input.Password = string(hashedPassword)
	}

	createdUser := validations.CreateUserInput{
		Fullname:  input.Fullname,
		Email:     input.Email,
		BirthDate: input.BirthDate,
		Password:  input.Password,
	}
	_, err := UserModel.CreateUser(&createdUser)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}

func (auth AuthController) Logout(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if err := UserModel.RemoveToken(&validations.LogoutInput{Userid: userID.(primitive.ObjectID)}); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}
