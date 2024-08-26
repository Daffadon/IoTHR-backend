package migrations

import (
	"IoTHR-backend/db"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() {
	client := db.GetDB()
	collection := client.Database("heartrate").Collection("user")
	if collection == nil {
		log.Fatal("Failed to create collection")
	}
	database := client.Database("heartrate")
	bucketOpts := options.GridFSBucket().SetName("resampled_ecg_files")
	_, err := gridfs.NewBucket(database, bucketOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database and collection created successfully")			
}
