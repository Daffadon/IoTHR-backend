package migrations

import (
	"IoTHR-backend/db"
	"fmt"
	"log"
)

func Init() {
	client := db.GetDB()
	collection := client.Database("heartrate").Collection("user")
	if collection == nil {
		log.Fatal("Failed to create collection")
	}
	fmt.Println("Database and collection created successfully")
}
