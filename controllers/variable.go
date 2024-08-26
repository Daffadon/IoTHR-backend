package controllers

import (
	"IoTHR-backend/errors"
	"IoTHR-backend/models"
)

var UserModel = new(models.User)
var TopicModel = new(models.Topic)
var PredictionModel = new(models.Prediction)
var ECGModel = new(models.ECG)
var ecgController = new(ECGController)
var topicController = new(TopicController)
var errorInstance = new(errors.ErrorInstance)
