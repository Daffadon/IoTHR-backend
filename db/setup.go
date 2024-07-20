package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func Init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to connect to database!")
	} else {
		fmt.Println("Connected to mongoDB!!!")
	}
	db = client
}

func GetUserCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("user")
}
func GetTopicCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("topic")
}

func GetDB() *mongo.Client {
	return db
}
