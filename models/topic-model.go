package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topic struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Date       time.Time          `bson:"date" json:"date"`
	RecordTime string             `bson:"recordTime,omitempty" json:"recordTime,omitempty"`
	EcgPlot    []float64          `bson:"ecgplot,omitempty" json:"ecgplot,omitempty"`
}

func (t Topic) CreateTopic(input *validations.CreateTopicWithUSerIDInput) (*Topic, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": input.Name}
	var existingUser User
	err := topicCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("Topic already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	topic := Topic{
		UserID: input.UserID,
		Name:   input.Name,
		Date:   time.Now(),
	}
	result, err := topicCollection.InsertOne(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to insert topic: %v", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to get inserted ID")
	}
	topic.ID = insertedID
	return &topic, nil
}

func (t Topic) GetTopicById(input *validations.GetTopicByIdInput) (*Topic, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topic Topic
	filter := bson.M{"_id": input.TopicID, "userId": input.UserID}
	err := topicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topic with ID %v: %v", input.TopicID, err)
	}
	return &topic, nil
}

func (t Topic) GetTopicList(input *validations.GetHistoryInput) ([]validations.HistoryReturn, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topicObject []validations.HistoryReturn
	for _, topicID := range input.TopicList {
		var topic Topic
		err := topicCollection.FindOne(ctx, bson.M{"_id": topicID}).Decode(&topic)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve topic with ID %v: %v", topicID, err)
		}
		toInsert := validations.HistoryReturn{
			TopicId:   topic.ID,
			TopicName: topic.Name,
		}
		topicObject = append(topicObject, toInsert)
	}
	return topicObject, nil
}

func (t Topic) UpdateECGdata(input *validations.UpdateECGInput) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errChannel := make(chan error)
	go func() {
		defer close(errChannel)
		filter := bson.M{"_id": input.TopicID, "userId": input.UserID}
		update := bson.M{
			"$push": bson.M{
				"ecgplot": bson.M{
					"$each": input.ECGPlot,
				},
			},
		}
		_, err := topicCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			errChannel <- fmt.Errorf("failed to update ECG data: %v", err)
			return
		}
		errChannel <- nil
	}()
	return <-errChannel
}

func (t Topic) UpdateRecordTime(input *validations.UpdateRecordTimeVal) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.TopicID, "userId": input.UserID}
	_, err := topicCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"recordTime": input.RecordTime}})
	if err != nil {
		return fmt.Errorf("failed to update record time: %v", err)
	}
	return nil
}

func (t Topic) ECGPrediction(input *validations.ECGPredictionInput) (*Prediction, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.TopicID, "userId": input.UserID}

	var topic Topic
	err := topicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topic with ID %v: %v", input.TopicID, err)
	}

	packet := validations.Payload{
		ECGPlot: topic.EcgPlot,
		Feature: input.Feature,
	}

	jsondata, err := json.Marshal(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	resp, err := http.Post(os.Getenv("MODEL_URL"), "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var predictionResult struct {
		N             float64       `json:"N"`
		S             float64       `json:"S"`
		V             float64       `json:"V"`
		F             float64       `json:"F"`
		Q             float64       `json:"Q"`
		Feature       string        `json:"feature"`
		SampleList    [][]float64   `json:"sample_list"`
		Annotation    [][][]float64 `json:"annotation"`
		InferenceTime float64       `json:"inferenceTime"`
		ModelVersion  string        `json:"modelVersion"`
	}

	err = json.Unmarshal(body, &predictionResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction result: %v", err)
	}
	createPrediction := &Prediction{
		TopicID:       input.TopicID,
		InferenceTime: predictionResult.InferenceTime,
		Prediction: struct {
			N float64 `bson:"N" json:"N"`
			S float64 `bson:"S" json:"S"`
			V float64 `bson:"V" json:"V"`
			F float64 `bson:"F" json:"F"`
			Q float64 `bson:"Q" json:"Q"`
		}{
			N: predictionResult.N,
			S: predictionResult.S,
			V: predictionResult.V,
			F: predictionResult.F,
			Q: predictionResult.Q,
		},
		Feature: predictionResult.Feature,
		Annotation: struct {
			N [][]float64 `bson:"N" json:"N"`
			S [][]float64 `bson:"S" json:"S"`
			V [][]float64 `bson:"V" json:"V"`
			F [][]float64 `bson:"F" json:"F"`
			Q [][]float64 `bson:"Q" json:"Q"`
		}{
			N: predictionResult.Annotation[0],
			S: predictionResult.Annotation[1],
			V: predictionResult.Annotation[2],
			F: predictionResult.Annotation[3],
			Q: predictionResult.Annotation[4],
		},
		SamplePlot: struct {
			N []float64 `bson:"N" json:"N"`
			S []float64 `bson:"S" json:"S"`
			V []float64 `bson:"V" json:"V"`
			F []float64 `bson:"F" json:"F"`
			Q []float64 `bson:"Q" json:"Q"`
		}{
			N: predictionResult.SampleList[0],
			S: predictionResult.SampleList[1],
			V: predictionResult.SampleList[2],
			F: predictionResult.SampleList[3],
			Q: predictionResult.SampleList[4],
		},
		ModelVersion: predictionResult.ModelVersion,
	}
	return createPrediction, nil
}
