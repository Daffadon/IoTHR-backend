package controllers

import (
	"IoTHR-backend/models"
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

var UserModel = new(models.User)

func (auth AuthController) Login(ctx *gin.Context) {
	var input validations.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	user, err := UserModel.GetUser(&input)
	if err == nil {
		errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if errCompare != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
			ctx.Abort()
			return
		}
		token, errJWT := utils.GenerateToken(user.ID, user.Role)
		if errJWT != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token!"})
			ctx.Abort()
			return
		}
		errUpdate := UserModel.UpdateTokenToUser(user.ID, token)
		if errUpdate != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token!"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"jwt": token})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
	ctx.Abort()
}

func (auth AuthController) Register(ctx *gin.Context) {
	var input validations.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if input.Password != input.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password and Confirm Password not match!"})
		ctx.Abort()
		return
	}

	if hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), 14); errHash != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password!"})
		ctx.Abort()
		return
	} else {
		input.Password = string(hashedPassword)
	}

	createdUser := validations.CreateUserInput{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: input.Password,
	}
	_, err := UserModel.CreateUser(&createdUser)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	ctx.Abort()
}

func (auth AuthController) Logout(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if err := UserModel.RemoveToken(&validations.LogoutInput{Userid: userID.(primitive.ObjectID)}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	ctx.Abort()
}
