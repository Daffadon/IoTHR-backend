package validations

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPredictionByTopicIdInput struct {
	TopicID primitive.ObjectID `json:"topicId" binding:"required"`
}

type PredictionListReturn struct {
	PredictionId primitive.ObjectID `json:"predictionId"`
	Feature      string             `json:"feature"`
}
