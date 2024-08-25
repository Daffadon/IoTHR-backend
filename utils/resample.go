package utils

import (
	"IoTHR-backend/validations"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ResampleECG(ecgplot []float64) ([]float64, error) {

	packet := &validations.ResampleECGDataInput{
		ECG_Plot: ecgplot,
	}
	jsondata, err := json.Marshal(packet)

	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error marshalling data")
	}
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}
	resp, err := http.Post(os.Getenv("RESAMPLE_URL"), "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error sending request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error reading response body")
	}
	var ecgPlot1D struct {
		ECG_Plot []float64 `json:"ecg_plot"`
	}
	err = json.Unmarshal(body, &ecgPlot1D)

	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error unmarshalling data")
	}
	return ecgPlot1D.ECG_Plot, nil
}
