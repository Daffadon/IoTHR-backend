package validations

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	Email    string `json:"email" `
	Fullname string `json:"fullname"`
}
type GetHistoryInput struct {
	TopicList []primitive.ObjectID `json:"topicList"`
}

type HistoryReturn struct {
	TopicId   primitive.ObjectID `json:"topicId"`
	TopicName string             `json:"topicName"`
}
