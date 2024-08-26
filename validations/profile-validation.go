package validations

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	Email     string `json:"email" `
	Fullname  string `json:"fullname"`
	BirthDate string `json:"birthDate"`
}
type GetHistoryInput struct {
	TopicList []primitive.ObjectID `json:"topicList"`
}

type HistoryReturn struct {
	TopicId   primitive.ObjectID `json:"topicId"`
	TopicName string             `json:"topicName"`
	Analyzed  bool               `json:"analyzed"`
}

type UserHistoryReturn struct {
	TopicId    primitive.ObjectID `json:"topicId"`
	TopicName  string             `json:"topicName"`
	Date       string             `json:"date"`
	RecordTime string             `json:"recordTime"`
	Analyzed   bool               `json:"analyzed"`
}
