package models

import (
	"IoTHR-backend/db"
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Date       string             `bson:"date" json:"date"`
	RecordTime string             `bson:"recordTime,omitempty" json:"recordTime,omitempty"`
	ECGFileId  primitive.ObjectID `bson:"ecgfileId,omitempty" json:"ecgfileId,omitempty"`
	Analyzed   bool               `bson:"analyzed" json:"analyzed"`
	Analysis   []struct {
		DoctorId   primitive.ObjectID `bson:"doctorId" json:"doctorId"`
		DoctorName string             `bson:"doctorName" json:"doctorName"`
		Comment    []string           `bson:"comment" json:"comment"`
	} `bson:"analysis,omitempty" json:"analysis,omitempty"`
}

func (t Topic) CreateTopic(input *validations.CreateTopicWithUSerIDInput) (*Topic, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": input.Name}
	var existingUser User
	err := topicCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, errorInstance.ReturnError(http.StatusConflict, "Topic already exists")
	}

	topic := Topic{
		UserID:   input.UserID,
		Name:     input.Name,
		Date:     time.Now().UTC().Local().Format("Monday, 02-01-2006 15:04:05 WIB"),
		Analyzed: false,
	}

	result, err := topicCollection.InsertOne(ctx, topic)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error inserting topic")
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Failed to get inserted ID")
	}
	topic.ID = insertedID
	return &topic, nil
}

func (t Topic) GetTopicById(input *validations.GetTopicByIdInput) (*validations.GetTopicByIDReturn, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topic validations.TopicVal
	filter := bson.M{"_id": input.TopicID, "userId": input.UserID}
	err := topicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "Topic not found")
	}
	if topic.ECGFileId.IsZero() {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "Recorded File Not Found	")
	}

	ecgPlots, err := utils.GetECGFileById(topic.ECGFileId)
	if err != nil {
		return nil, err
	}
	returnedData := &validations.GetTopicByIDReturn{
		Topic:   topic,
		ECGPlot: ecgPlots,
	}
	return returnedData, nil
}

func (t Topic) GetTopicByIdForDoctor(input *primitive.ObjectID) (*validations.GetTopicByIDReturn, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topic validations.TopicVal
	filter := bson.M{"_id": input}
	err := topicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusNotFound, "Topic not found")
	}
	ecgPlots, err := utils.GetECGFileById(topic.ECGFileId)
	if err != nil {
		return nil, err
	}

	returnedData := &validations.GetTopicByIDReturn{
		Topic:   topic,
		ECGPlot: ecgPlots,
	}

	return returnedData, nil
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
			return nil, errorInstance.ReturnError(http.StatusNotFound, "Topic not found")
		}
		toInsert := validations.HistoryReturn{
			TopicId:   topic.ID,
			TopicName: topic.Name,
			Analyzed:  topic.Analyzed,
		}
		topicObject = append(topicObject, toInsert)
	}
	return topicObject, nil
}

func (t Topic) GetUserTopicList(input []primitive.ObjectID) ([]validations.UserHistoryReturn, error) {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topicObject []validations.UserHistoryReturn
	for _, topicID := range input {
		var topic Topic
		err := topicCollection.FindOne(ctx, bson.M{"_id": topicID}).Decode(&topic)
		if err != nil {
			return nil, errorInstance.ReturnError(http.StatusNotFound, "Topic not found")
		}
		toInsert := validations.UserHistoryReturn{
			TopicId:    topic.ID,
			TopicName:  topic.Name,
			Date:       topic.Date,
			RecordTime: topic.RecordTime,
			Analyzed:   topic.Analyzed,
		}
		topicObject = append(topicObject, toInsert)
	}
	return topicObject, nil

}

func (t Topic) UpdateTopicAnalyzeData(input *validations.UpdateAnalyzeData) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"analyzed": input.Analyzed,
		},
	}
	filter := bson.M{"_id": input.TopicID}
	_, err := topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating topic")
	}
	return nil
}

func (t Topic) UpdateTopicAnalyzeComment(input *validations.UpdateAnalyzeComment) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingAnalysis := bson.M{"analysis.doctorId": input.DoctorID}
	count, err := topicCollection.CountDocuments(ctx, existingAnalysis)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error checking if doctorId already exists")
	}
	if count > 0 {
		update := bson.M{
			"$push": bson.M{
				"analysis.$.comment": bson.M{
					"$each": input.Comment,
				},
			},
		}
		filter := bson.M{"_id": input.TopicID, "analysis.doctorId": input.DoctorID}
		_, err := topicCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating topic")
		}
		return nil
	}

	if count == 0 {
		newAnalysis := struct {
			DoctorID   primitive.ObjectID `bson:"doctorId" json:"doctorId"`
			DoctorName string             `bson:"doctorName" json:"doctorName"`
			Comment    []string           `bson:"comment" json:"comment"`
		}{
			DoctorID:   input.DoctorID,
			DoctorName: input.DoctorName,
			Comment:    input.Comment,
		}
		update := bson.M{
			"$push": bson.M{
				"analysis": newAnalysis,
			},
		}
		filter := bson.M{"_id": input.TopicID}
		_, err := topicCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating topic")
		}
	}
	return nil
}

func (t Topic) DeleteTopicAnalyzeComment(input *validations.DeleteAnalyzeComment) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.TopicID, "analysis.doctorId": input.DoctorID}
	update := bson.M{"$unset": bson.M{"analysis.$.comment." + fmt.Sprint(input.CommentIndex): ""}}
	_, err := topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error deleting comment")
	}

	err = utils.CleanupComment(input.DoctorID)
	if err != nil {
		return err
	}

	var count int64
	countFilter := bson.M{
		"analysis.doctorId": input.DoctorID,
		"analysis.comment":  bson.M{"$ne": []interface{}{}},
	}

	count, err = topicCollection.CountDocuments(ctx, countFilter)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error counting comments")
	}

	if count == 0 {
		removeFilter := bson.M{
			"analysis.doctorId": input.DoctorID,
		}
		_, err = topicCollection.UpdateMany(ctx, removeFilter, bson.M{"$pull": bson.M{"analysis": bson.M{"doctorId": input.DoctorID}}})
		if err != nil {
			return errorInstance.ReturnError(http.StatusInternalServerError, "Error removing doctor from analysis")
		}
	}
	return nil
}

func (t Topic) UpdateRecordTime(input *validations.UpdateRecordTimeVal) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.TopicID, "userId": input.UserID}
	_, err := topicCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"recordTime": input.RecordTime}})
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating record time")
	}
	return nil
}

func (t Topic) UpdateECGFileID(input *validations.UpdateECGFileID) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.TopicID}
	_, err := topicCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"ecgfileId": input.ECGFileID}})
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating ECG file ID")
	}
	return nil
}

func (t Topic) DeleteTopic(topicId primitive.ObjectID) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": topicId}
	_, err := topicCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error deleting topic")
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
		return nil, errorInstance.ReturnError(http.StatusNotFound, "Topic not found")
	}

	ecgPlot, err := utils.GetECGFileById(topic.ECGFileId)
	if err != nil {
		return nil, err
	}

	packet := validations.Payload{
		ECGPlot: ecgPlot,
		Feature: input.Feature,
	}

	jsondata, err := json.Marshal(packet)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error marshalling data")
	}

	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	resp, err := http.Post(os.Getenv("MODEL_URL"), "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error sending request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error reading response body")
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
		return nil, errorInstance.ReturnError(http.StatusInternalServerError, "Error unmarshalling data")
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
