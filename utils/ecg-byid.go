package utils

import (
	"IoTHR-backend/db"
	"bytes"
	"encoding/json"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetECGFileById(fileID primitive.ObjectID) ([]float64, error) {

	bucket := db.GetECGFileBucket()
	var buf bytes.Buffer
	downloadStream, err := bucket.OpenDownloadStream(fileID)
	if err != nil {
		log.Fatal("Error opening download stream from GridFS:", err)
	}
	defer downloadStream.Close()

	_, err = io.Copy(&buf, downloadStream)
	if err != nil {
		log.Fatal("Error reading from GridFS download stream:", err)
	}

	var ecgPlot struct {
		ECG_Plot []float64 `json:"ecg_plot" bson:"ecg_plot"`
	}
	err = json.Unmarshal(buf.Bytes(), &ecgPlot)
	if err != nil {
		log.Fatal("Error unmarshalling JSON data:", err)
	}

	return ecgPlot.ECG_Plot, nil
}
