package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/validations"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topic struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Date       time.Time          `bson:"date" json:"date"`
	Prediction struct {
		N float64 `bson:"N" json:"N"`
		S float64 `bson:"S" json:"S"`
		V float64 `bson:"V" json:"V"`
		F float64 `bson:"F" json:"F"`
		Q float64 `bson:"Q" json:"Q"`
	} `bson:"prediction,omitempty" json:"prediction,omitempty"`
	Feature    string    `bson:"feature,omitempty" json:"feature,omitempty"`
	Runtime    float64   `bson:"runtime,omitempty" json:"runtime,omitempty"`
	EcgPlot    []float64 `bson:"ecgplot,omitempty" json:"ecgplot,omitempty"`
	Annotation struct {
		N [][]float64 `bson:"N" json:"N"`
		S [][]float64 `bson:"S" json:"S"`
		V [][]float64 `bson:"V" json:"V"`
		F [][]float64 `bson:"F" json:"F"`
		Q [][]float64 `bson:"Q" json:"Q"`
	} `bson:"annotation,omitempty										" json:"annotation,omitempty"`
	SamplePlot struct {
		N []interface{} `bson:"N" json:"N"`
		S []interface{} `bson:"S" json:"S"`
		V []interface{} `bson:"V" json:"V"`
		F []interface{} `bson:"F" json:"F"`
		Q []interface{} `bson:"Q" json:"Q"`
	} `bson:"sample_plot,omitempty" json:"sample_plot,omitempty"`
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
