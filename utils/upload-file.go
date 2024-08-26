package utils

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadFile(TopicId primitive.ObjectID, ecgPlot1D []float64) (primitive.ObjectID, error) {

	ecgJson := validations.ResampleECGDataInput{
		ECG_Plot: ecgPlot1D,
	}

	jsonData, err := json.Marshal(ecgJson)
	if err != nil {
		return primitive.NilObjectID, errorInstance.ReturnError(http.StatusInternalServerError, "Error marshalling data")
	}

	bucket := db.GetECGFileBucket()
	uploadStream, err := bucket.OpenUploadStream(TopicId.Hex() + ".json")
	if err != nil {
		return primitive.NilObjectID, errorInstance.ReturnError(http.StatusInternalServerError, "Error opening upload stream")
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(jsonData)
	if err != nil {
		return primitive.NilObjectID, errorInstance.ReturnError(http.StatusInternalServerError, "Error writing JSON data to GridFS")
	}

	fileID, ok := uploadStream.FileID.(primitive.ObjectID)
	if !ok {
		return primitive.NewObjectID(), errorInstance.ReturnError(http.StatusInternalServerError, "Failed to get file ID")
	}
	fmt.Println("File uploaded successfully with ID:", fileID.Hex())
	return fileID, nil
}
