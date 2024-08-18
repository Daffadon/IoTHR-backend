package validations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicVal struct {
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

type CreateTopicInput struct {
	Name string `json:"name" bson:"name"`
}
type CreateTopicWithUSerIDInput struct {
	Name   string             `json:"name" bson:"name"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
}

type InsertECGDataInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
}

type UpdateECGInput struct {
	TopicID  primitive.ObjectID `json:"topicId" bson:"topicId"`
	ECGPlot  []float64          `json:"ECGplot" bson:"ECGplot"`
	Sequence int                `json:"sequence" bson:"sequence"`
}

type GetTopicByIdInput struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
	UserID  primitive.ObjectID `json:"userId" bson:"userId"`
}

type Prediction struct {
	TopicID primitive.ObjectID `json:"topicId" bson:"topicId"`
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
	CommentIndex int `json:"commentIndex" bson:"commentIndex"`
}

type DeleteAnalyzeComment struct {
	TopicID      primitive.ObjectID `json:"topicId" bson:"topicId"`
	DoctorID     primitive.ObjectID `json:"doctorId" bson:"doctorId"`
	CommentIndex int                `json:"commentIndex" bson:"commentIndex"`
}

type UpdateECGFileID struct {
	TopicID   primitive.ObjectID `json:"topicId" bson:"topicId"`
	ECGFileID primitive.ObjectID `json:"ecgFileId" bson:"ecgFileId"`
}

type GetTopicByIDReturn struct {
	Topic   TopicVal  `json:"topic" bson:"topic"`
	ECGPlot []float64 `json:"ecg_plot" bson:"ecg_plot"`
}
