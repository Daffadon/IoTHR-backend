package validations

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoginInput struct {
	Email    string `json:"email" binding:"required" type:"email"`
	Password string `json:"password" binding:"required"`
}
type RegisterInput struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required" type:"email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}
type CreateUserInput struct {
	Email    string `json:"email" binding:"required" type:"email"`
	Fullname string `json:"fullname" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LogoutInput struct {
	Userid primitive.ObjectID `json:"user_id" binding:"required"`
}
