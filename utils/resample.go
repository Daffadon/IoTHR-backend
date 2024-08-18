package utils

import (
	"IoTHR-backend/validations"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	resp, err := http.Post(os.Getenv("RESAMPLE_URL"), "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var ecgPlot1D struct {
		ECG_Plot []float64 `json:"ecg_plot"`
	}
	err = json.Unmarshal(body, &ecgPlot1D)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction result: %v", err)
	}
	return ecgPlot1D.ECG_Plot, nil
}
