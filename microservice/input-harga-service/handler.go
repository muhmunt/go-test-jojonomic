package main

import (
	"encoding/json"
	"fmt"
	"input-harga-service/helper"
	"input-harga-service/model"
	"input-harga-service/request"
	"input-harga-service/service"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type priceHandler struct {
	priceService service.Service
}

func NewPrice(service service.Service) *priceHandler {
	return &priceHandler{service}
}

func (h *priceHandler) handleInputPriceRequest(c *gin.Context) {
	var request request.CreatePriceRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	_, err = h.priceService.FindById(request.AdminID)

	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	priceData := model.Price{
		AdminID:      request.AdminID,
		HargaTopup:   request.HargaTopup,
		HargaBuyback: request.HargaBuyback,
	}

	bytes, err := json.Marshal(priceData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal JSON"})
		return
	}

	producer, err := CreateSyncProducer([]string{"kafka:9092"})
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Kafka producer"})
		return
	}
	defer producer.Close()

	err = SendMessageToKafka(producer, "input-harga", request.AdminID, bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message to Kafka"})
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	mu.Lock()
	responseChannels[request.AdminID] = responseCh
	mu.Unlock()

	select {
	case responseMsg := <-responseCh:
		fmt.Println(responseMsg.Value)
		if string(responseMsg.Value) == "false" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"reff_id": helper.GenShortId(),
				"message": "Kafka not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"reff_id": helper.GenShortId(),
		})

	case <-time.After(10 * time.Second):
		mu.Lock()
		delete(responseChannels, request.AdminID)
		mu.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "timeout waiting for response"})
	}
}
