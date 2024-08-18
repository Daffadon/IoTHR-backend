package controllers

import (
	"IoTHR-backend/utils"
	"IoTHR-backend/validations"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	TopicId  primitive.ObjectID `json:"topicId"`
	Type     string             `json:"type"`
	ECGPlot  []float64          `json:"ECGPlot"`
	JWT      string             `json:"jwt"`
	Sequence int                `json:"sequence"`
}

type WebsocketController struct{}

func (w WebsocketController) UpdateECGPlot(c *gin.Context) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}
	ctx := c.Request.Context()

	conn, err := websocket.Accept(c.Writer, c.Request, opts)
	if err != nil {
		log.Println("Error accepting WebSocket connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept WebSocket connection"})
		return
	}

	defer conn.Close(websocket.StatusInternalError, "Connection closed with error")

	log.Println("WebSocket connection established")

	for {
		var msg Message

		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			log.Fatalln("Error reading message:", err)
			return
		}
		if msg.Type == "close" {
			claims, err := utils.ValidateJWT(msg.JWT)
			if err != nil {
				log.Println("Error validating JWT:", err)
				return
			}
			ecgController.ResampleECGData(msg.TopicId)
			topicController.PredictionForWS(msg.TopicId, claims.UserId)
			return
		}

		dataToUpdate := &validations.UpdateECGInput{
			TopicID:  msg.TopicId,
			ECGPlot:  msg.ECGPlot,
			Sequence: msg.Sequence,
		}

		sequence, err := ECGModel.UpdateECGdata(dataToUpdate)

		if err != nil {
			fmt.Print("Error updating ECG data:", err)
			return
		}

		err = wsjson.Write(ctx, conn, sequence)
		if err != nil {
			log.Fatalln("Error sending message:", err)
			return
		}
	}
}
