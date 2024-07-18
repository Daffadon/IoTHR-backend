package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"fmt"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique;not null" `
	Fullname string `json:"fullname" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'user';not null"`
	Token    string `json:"token" gorm:"default:null"`
}

func (u User) GetUser(input *validations.LoginInput) (*User, error) {
	db := db.GetDB()
	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) GetUserByID(id uint) (*User, error) {
	db := db.GetDB()
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) UpdateTokenToUser(id uint, token string) error {
	db := db.GetDB()
	if err := db.Model(&User{}).Where("id = ?", id).Update("token", token); err != nil {
		return err.Error
	}
	return nil
}

func (u User) CreateUser(input *validations.CreateUserInput) (*User, error) {
	db := db.GetDB()
	user := User{Fullname: input.Fullname, Email: input.Email, Password: input.Password}
	if err := db.Create(&user).Error; err != nil {
		fmt.Println(err.Error())
		if utils.IsUniqueConstraintError(err) {
			return nil, fmt.Errorf("email already exists")
		}
	}
	return &user, nil
}
func (u User) RemoveToken(input *validations.LogoutInput) error {
	db := db.GetDB()
	if err := db.Model(&User{}).Where("id = ?", input.Userid).Update("token", nil); err != nil {
		return err.Error
	}
	return nil
}

func IsLoggedOut(token string) bool {
	var count int64
	db := db.GetDB()
	db.Model(&User{}).Where("token = ?", token).Count(&count)
	return count > 0
}
