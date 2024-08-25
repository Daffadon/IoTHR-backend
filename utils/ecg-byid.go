package utils

import (
	"IoTHR-backend/db"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetECGFileById(fileID primitive.ObjectID) ([]float64, error) {

	bucket := db.GetECGFileBucket()
	var buf bytes.Buffer
	downloadStream, err := bucket.OpenDownloadStream(fileID)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error opening download stream from GridFS")
	}
	defer downloadStream.Close()

	_, err = io.Copy(&buf, downloadStream)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error reading from GridFS download stream")
	}

	var ecgPlot struct {
		ECG_Plot []float64 `json:"ecg_plot" bson:"ecg_plot"`
	}
	err = json.Unmarshal(buf.Bytes(), &ecgPlot)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error unmarshalling JSON data")
	}

	return ecgPlot.ECG_Plot, nil
}
