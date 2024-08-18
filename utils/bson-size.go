package utils

import "go.mongodb.org/mongo-driver/bson"

func BsonSize(value interface{}) int {
	data, _ := bson.Marshal(value)
	return len(data)
}
