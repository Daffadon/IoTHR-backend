package utils

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadFile(TopicId primitive.ObjectID, ecgPlot1D []float64) (primitive.ObjectID, error) {

	ecgJson := validations.ResampleECGDataInput{
		ECG_Plot: ecgPlot1D,
	}

	jsonData, err := json.Marshal(ecgJson)
	if err != nil {
		log.Fatal("Error marshalling ECG data to JSON:", err)
		return primitive.NilObjectID, err
	}

	bucket := db.GetECGFileBucket()
	uploadStream, err := bucket.OpenUploadStream(TopicId.Hex() + ".json")
	if err != nil {
		log.Fatal("Error opening upload stream to GridFS:", err)
		return primitive.NilObjectID, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(jsonData)

	if err != nil {
		log.Fatal("Error writing JSON data to GridFS:", err)
		return primitive.NilObjectID, err
	}

	fileID, ok := uploadStream.FileID.(primitive.ObjectID)
	if !ok {
		log.Fatal("Error: FileID is not of type primitive.ObjectID")
		return primitive.NilObjectID, err
	}
	fmt.Println("File uploaded successfully with ID:", fileID.Hex())

	return fileID, nil
}
