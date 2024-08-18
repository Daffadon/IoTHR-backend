package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ECG struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TopicID  primitive.ObjectID `json:"topicId,omitempty" bson:"topicId,omitempty"`
	ECG_Plot []float64          `json:"ECGPlot,omitempty" bson:"ECGPlot,omitempty"`
	Sequence int                `json:"sequence,omitempty" bson:"sequence,omitempty"`
}

func (e ECG) CreateECGData(input *validations.InsertECGDataInput) error {
	ecgColletion := db.GetECGCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ecgData := ECG{
		TopicID:  input.TopicID,
		Sequence: 1,
	}
	_, err := ecgColletion.InsertOne(ctx, ecgData)
	if err != nil {
		return fmt.Errorf("failed to insert ECG data: %v", err)
	}
	return nil
}

func (e ECG) UpdateECGdata(input *validations.UpdateECGInput) (int, error) {
	ecgCollection := db.GetECGCollection()
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	// defer cancel()

	filter := bson.M{"topicId": input.TopicID, "sequence": input.Sequence}
	var result bson.M

	err := ecgCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("no document found with the given TopicID and sequence")
		}
		return 0, fmt.Errorf("failed to fetch document: %v", err)
	}
	approxSize := utils.BsonSize(result)
	sequenceID := input.Sequence
	if seq, ok := result["sequence"].(int); ok {
		sequenceID = seq
	}

	if approxSize+len(input.ECGPlot) > 15*1024*1024 {
		newSequenceID := sequenceID + 1

		newDoc := bson.M{
			"topicId":  input.TopicID,
			"ecgplot":  input.ECGPlot,
			"sequence": newSequenceID,
		}
		fmt.Println("Create Document")

		_, err := ecgCollection.InsertOne(ctx, newDoc)
		if err != nil {
			return 0, fmt.Errorf("failed to create new ECG data document: %v", err)

		}
		return sequenceID, nil
	} else {
		update := bson.M{
			"$push": bson.M{
				"ecgplot": bson.M{
					"$each": input.ECGPlot,
				},
			},
		}
		_, err = ecgCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return 0, fmt.Errorf("failed to update ECG data: %v", err)
		}
		return sequenceID, nil
	}
}

func (e ECG) ECGMergePlot(TopicId primitive.ObjectID) ([]float64, error) {
	ecgCollection := db.GetECGCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"topicId": TopicId,
			},
		},
		{
			"$project": bson.M{
				"ecgplot": 1,
			},
		},
	}

	cursor, err := ecgCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate ECG data: %v", err)
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode ECG data: %v", err)
	}

	ecgPlots := results[0]
	rawEcgPlotData, exists := ecgPlots["ecgplot"]
	if !exists {
		fmt.Println("Error: ecgplot data not found.")
		return nil, fmt.Errorf("ecgplot data not found")
	}

	// Attempt to cast the raw data to []interface{}
	ecgPlotData, ok := rawEcgPlotData.(primitive.A)
	if !ok {
		fmt.Println("Error: ecgplot data is not in the expected format.")
		return nil, fmt.Errorf("invalid ecgplot format")
	}

	// Convert []interface{} to []float64 or []int
	var ecgPlotFloats []float64
	for _, v := range ecgPlotData {
		switch val := v.(type) {
		case int:
			ecgPlotFloats = append(ecgPlotFloats, float64(val))
		case float64:
			ecgPlotFloats = append(ecgPlotFloats, val)
		default:
			fmt.Println("Error: ecgplot contains unsupported value types.")
			return nil, fmt.Errorf("ecgplot contains unsupported value types")
		}
	}

	return ecgPlotFloats, nil

}
