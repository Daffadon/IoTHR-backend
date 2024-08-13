package validations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type Prediction struct {
	TopicID    primitive.ObjectID `json:"topicId" bson:"topicId"`
	RecordTime string             `json:"recordTime" bson:"recordTime"`
	Feature    string             `json:"feature" bson:"feature"`
}
type UpdateRecordTimeVal struct {
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	TopicID    primitive.ObjectID `json:"topicId" bson:"topicId"`
	RecordTime string             `json:"recordTime" bson:"recordTime"`
}
type ECGPredictionInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
	Feature string             `json:"feature" bson:"feature"`
}

type Payload struct {
	ECGPlot []float64 `json:"ecg_plot" bson:"ecg_plot"`
	Feature string    `json:"feature" bson:"feature"`
}

type UpdateAnalyzeData struct {
	TopicID  primitive.ObjectID `json:"topicId" bson:"topicId"`
	Analyzed bool               `json:"analyzed" bson:"analyzed"`
}

type UpdateAnalyzeCommentInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	Comment string             `json:"comment" bson:"comment"`
}
type UpdateAnalyzeComment struct {
	TopicID    primitive.ObjectID `json:"topicId" bson:"topicId"`
	DoctorID   primitive.ObjectID `json:"doctorId" bson:"doctorId"`
	DoctorName string             `json:"doctorName" bson:"doctorName"`
	Comment    []string           `json:"comment" bson:"comment"`
}

type DeleteAnalyzeCommentInput struct {
	// TopicID      primitive.ObjectID `json:"topicId" bson:"topicId"`
	CommentIndex int `json:"commentIndex" bson:"commentIndex"`
}

type DeleteAnalyzeComment struct {
	TopicID      primitive.ObjectID `json:"topicId" bson:"topicId"`
	DoctorID     primitive.ObjectID `json:"doctorId" bson:"doctorId"`
	CommentIndex int                `json:"commentIndex" bson:"commentIndex"`
}
