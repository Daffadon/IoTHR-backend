package controllers

import (
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ECGController struct{}

func (e ECGController) CreateECGData(TopicId *primitive.ObjectID) error {
	ecgData := &validations.InsertECGDataInput{
		TopicID: *TopicId,
	}
	err := ECGModel.CreateECGData(ecgData)
	if err != nil {
		return err
	}
	return nil
}

func (e ECGController) ResampleECGData(TopicId primitive.ObjectID) error {

	ecgPlots, err := ECGModel.ECGMergePlot(TopicId)
	if err != nil {
		return err
	}
	ecgPlot1D, err := utils.ResampleECG(ecgPlots)
	if err != nil {
		return err
	}

	fileId, err := utils.UploadFile(TopicId, ecgPlot1D)
	if err != nil {
		return err
	}
	topicToUpdate := &validations.UpdateECGFileID{
		TopicID:   TopicId,
		ECGFileID: fileId,
	}
	err = TopicModel.UpdateECGFileID(topicToUpdate)
	if err != nil {
		return err
	}
	return nil
}
