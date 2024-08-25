package utils

import (
	"IoTHR-backend/db"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CleanupComment(doctorId primitive.ObjectID) error {
	topicCollection := db.GetTopicCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc struct {
		Analysis []struct {
			DoctorId   primitive.ObjectID `bson:"doctorId" json:"doctorId"`
			DoctorName string             `bson:"doctorName" json:"doctorName"`
			Comment    []string           `bson:"comment" json:"comment"`
		} `bson:"analysis,omitempty" json:"analysis,omitempty"`
	}

	filter := bson.M{"analysis.doctorId": doctorId}
	err := topicCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return errorInstance.ReturnError(http.StatusNotFound, "Doctor not found")
	}

	updatedAnalysis := make([]struct {
		DoctorId   primitive.ObjectID `bson:"doctorId" json:"doctorId"`
		DoctorName string             `bson:"doctorName" json:"doctorName"`
		Comment    []string           `bson:"comment" json:"comment"`
	}, 0)

	for _, analysis := range doc.Analysis {
		if analysis.DoctorId == doctorId {
			cleanedComment := make([]string, 0)
			for _, c := range analysis.Comment {
				if c != "" {
					cleanedComment = append(cleanedComment, c)
				}
			}

			updatedAnalysis = append(updatedAnalysis, struct {
				DoctorId   primitive.ObjectID `bson:"doctorId" json:"doctorId"`
				DoctorName string             `bson:"doctorName" json:"doctorName"`
				Comment    []string           `bson:"comment" json:"comment"`
			}{
				DoctorId:   analysis.DoctorId,
				DoctorName: analysis.DoctorName,
				Comment:    cleanedComment,
			})
		} else {
			updatedAnalysis = append(updatedAnalysis, analysis)
		}
	}

	update := bson.M{"$set": bson.M{"analysis": updatedAnalysis}}

	_, err = topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorInstance.ReturnError(http.StatusInternalServerError, "Error updating comment")
	}
	return nil
}
