package validations

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateTopicInput struct {
	Name string `json:"name" bson:"name"`
}
type CreateTopicWithUSerIDInput struct {
	Name   string             `json:"name" bson:"name"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
}

type InsertECGDataInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
	ECGPlot []float64          `json:"ECGplot" bson:"ECGplot"`
}

type UpdateECGInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
	ECGPlot []float64          `json:"ECGplot" bson:"ECGplot"`
}

type GetTopicByIdInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
}

type TopicId struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
}
type ECGPredictionInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
}

type Payload struct {
	ECGPlot []float64 `json:"ecg_plot" bson:"ecg_plot"`
}
