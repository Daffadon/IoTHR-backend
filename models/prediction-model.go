package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prediction struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TopicID       primitive.ObjectID `json:"topicId,omitempty" bson:"topicId,omitempty"`
	InferenceTime float64            `bson:"inferenceTime,omitempty" json:"inferenceTime,omitempty"`
	Prediction    struct {
		N float64 `bson:"N" json:"N"`
		S float64 `bson:"S" json:"S"`
		V float64 `bson:"V" json:"V"`
		F float64 `bson:"F" json:"F"`
		Q float64 `bson:"Q" json:"Q"`
	} `bson:"prediction,omitempty" json:"prediction,omitempty"`
	Feature    string `bson:"feature,omitempty" json:"feature,omitempty"`
	Annotation struct {
		N [][]float64 `bson:"N" json:"N"`
		S [][]float64 `bson:"S" json:"S"`
		V [][]float64 `bson:"V" json:"V"`
		F [][]float64 `bson:"F" json:"F"`
		Q [][]float64 `bson:"Q" json:"Q"`
	} `bson:"annotation,omitempty" json:"annotation,omitempty"`
	SamplePlot struct {
		N []float64 `bson:"N" json:"N"`
		S []float64 `bson:"S" json:"S"`
		V []float64 `bson:"V" json:"V"`
		F []float64 `bson:"F" json:"F"`
		Q []float64 `bson:"Q" json:"Q"`
	} `bson:"sample_plot,omitempty" json:"sample_plot,omitempty"`
	ModelVersion string `bson:"modelVersion,omitempty" json:"modelVersion,omitempty"`
}

func (p Prediction) CreatePrediction(input *Prediction) (*primitive.ObjectID, error) {
	predictionCollection := db.GetPredictionCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := predictionCollection.InsertOne(ctx, input)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error inserting prediction")
	}

	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (p Prediction) GetPredictionByTopicId(input *primitive.ObjectID) (*[]validations.PredictionListReturn, error) {
	predictionCollection := db.GetPredictionCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"topicId": input}
	cursor, err := predictionCollection.Find(ctx, filter)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error fetching prediction")
	}
	var predictions []validations.PredictionListReturn
	for cursor.Next(ctx) {
		var prediction Prediction
		cursor.Decode(&prediction)
		predictions = append(predictions, validations.PredictionListReturn{
			PredictionId: prediction.ID,
			Feature:      prediction.Feature,
		})
	}
	return &predictions, nil
}

func (p Prediction) GetPredictionByPredictionId(input *primitive.ObjectID) (*Prediction, error) {
	predictionCollection := db.GetPredictionCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input}
	var prediction Prediction
	err := predictionCollection.FindOne(ctx, filter).Decode(&prediction)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "Prediction not found")
	}
	return &prediction, nil

}
