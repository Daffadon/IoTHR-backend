package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string               `json:"email" bson:"email"`
	Fullname string               `json:"fullname" bson:"fullname"`
	Password string               `json:"password" bson:"password"`
	Role     string               `json:"role" bson:"role"`
	Token    string               `json:"token,omitempty" bson:"token,omitempty"`
	TopicID  []primitive.ObjectID `json:"topicId,omitempty" bson:"topicId,omitempty"`
}

func (u User) CreateUser(input *validations.CreateUserInput) (*User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := User{Fullname: input.Fullname, Email: input.Email, Password: input.Password, Role: "user"}

	filter := bson.M{"email": input.Email}
	var existingUser User
	err := userCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	result, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, fmt.Errorf("failed to get inserted ID")
	}

	user.ID = insertedID
	return &user, nil
}

func (u User) GetUser(input *validations.LoginInput) (*User, error) {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	filter := bson.M{"email": input.Email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
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
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
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
		return fmt.Errorf("error updating user token: %v", err)
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
		return fmt.Errorf("error updating user token: %v", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("user not found or token already removed")
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

func (u User) UpdateTopicID(userid primitive.ObjectID, topicid primitive.ObjectID) error {
	userCollection := db.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": userid}
	update := bson.M{"$push": bson.M{"topicId": topicid}}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user topicId: %v", err)
	}
	return nil
}
	