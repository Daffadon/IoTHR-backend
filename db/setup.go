package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var db *mongo.Client

func Init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	MONGO_ADMIN_URI := os.Getenv("MONGO_ADMIN_URI")
	clientOptions := options.Client().ApplyURI(MONGO_ADMIN_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("Failed to connect to database!")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	adminDB := client.Database("heartrate")
	newUser := os.Getenv("APP_DB_USERNAME")
	password := os.Getenv("APP_DB_PASSWORD")

	checkUserCommand := bson.D{
		{Key: "usersInfo", Value: newUser},
	}

	var userResult bson.M
	if err := adminDB.RunCommand(context.Background(), checkUserCommand).Decode(&userResult); err != nil {
		log.Fatalf("Failed to check user: %v", err)
	}

	users, ok := userResult["users"].(bson.A)
	if ok && len(users) > 0 {
		log.Default().Println("User already exists!")
	} else if len(users) == 0 {
		result := adminDB.RunCommand(context.TODO(), bson.D{
			{Key: "createUser", Value: newUser},
			{Key: "pwd", Value: password},
			{Key: "roles", Value: bson.A{
				bson.D{{Key: "role", Value: "readWrite"}, {Key: "db", Value: os.Getenv("DATABASE")}},
				bson.D{{Key: "role", Value: "dbAdmin"}, {Key: "db", Value: os.Getenv("DATABASE")}},
			}},
		})
		if err = result.Err(); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Println("User created successfully!")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	hrClientOptions := options.Client().ApplyURI(MONGO_URI)
	hrClient, err := mongo.Connect(context.Background(), hrClientOptions)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db = hrClient
}

func GetUserCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("user")
}
func GetTopicCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("topic")
}
func GetPredictionCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("prediction")
}
func GetECGCollection() *mongo.Collection {
	return db.Database("heartrate").Collection("ecg")
}
func GetECGFileBucket() *gridfs.Bucket {
	database := db.Database("heartrate")
	bucketOpts := options.GridFSBucket().SetName("resampled_ecg_files")
	bucket, err := gridfs.NewBucket(database, bucketOpts)
	if err != nil {
		log.Fatal(err)
	}
	return bucket
}
func GetDB() *mongo.Client {
	return db
}
