package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/errors"
	"IoTHR-backend/validations"
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errorInstance = new(errors.ErrorInstance)

type User struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string               `json:"email" bson:"email"`
	Fullname  string               `json:"fullname" bson:"fullname"`
	BirthDate string               `json:"birthDate" bson:"birthDate"`
	Password  string               `json:"password" bson:"password"`
	Role      string               `json:"role" bson:"role"`
	Token     string               `json:"token,omitempty" bson:"token,omitempty"`
	TopicID   []primitive.ObjectID `json:"topicId,omitempty" bson:"topicId,omitempty"`
}

func (u User) CreateUser(input *validations.CreateUserInput) (*User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := User{Fullname: input.Fullname, Email: input.Email, BirthDate: input.BirthDate, Password: input.Password, Role: "user"}

	filter := bson.M{"email": input.Email}
	var existingUser User
	err := userCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, errorInstance.ReturnError(http.StatusConflict, "User already exists")
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error inserting user")
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Failed to get inserted ID")
	}
	user.ID = insertedID
	return &user, nil
}

func (u User) GetUsers() (*[]User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{
		"email":     1,
		"fullname":  1,
		"birthDate": 1,
		"_id":       1,
	}

	cursor, err := userCollection.Find(ctx, bson.M{"role": "user"}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error getting users")
	}

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error decoding users")
	}
	return &users, nil
}

func (u User) GetUser(input *validations.LoginInput) (*User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	filter := bson.M{"email": input.Email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "User Not Found")
	}
	return &user, nil
}

func (u User) GetUserByID(id primitive.ObjectID) (*User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	filter := bson.M{"_id": id}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "User Not Found")
	}
	return &user, nil
}

func (u User) UpdateTokenToUser(userid primitive.ObjectID, token string) error {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": userid}
	update := bson.M{"$set": bson.M{"token": token}}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating user token")
	}
	return nil
}

func (u User) RemoveToken(input *validations.LogoutInput) error {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.Userid}
	update := bson.M{"$unset": bson.M{"token": ""}}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating user token")
	}

	if result.ModifiedCount == 0 {
		return errorInstance.ReturnError(http.StatusNotFound, "User not found or token already removed")
	}
	return nil
}
func (u User) UpdateTopicID(userid primitive.ObjectID, topicid primitive.ObjectID) error {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": userid}
	update := bson.M{"$push": bson.M{"topicId": topicid}}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating user topicId")
	}
	return nil
}

func (u User) DeleteTopicID(topicId primitive.ObjectID) error {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"topicId": topicId}
	update := bson.M{"$pull": bson.M{"topicId": topicId}}

	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error deleting user topicId")
	}
	return nil
}

func (u User) IsLoggedOut(token string) bool {
	collection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"token": token}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error counting documents: %v\n", err)
		return false
	}
	return count > 0
}
