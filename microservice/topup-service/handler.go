package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"topup-service/helper"
	"topup-service/model"
	"topup-service/request"
	"topup-service/service"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type topupHandler struct {
	priceService       service.PriceService
	accountService     service.AccountService
	transactionService service.TransactionService
}

func NewTopup(priceService service.PriceService, accountService service.AccountService, transactionService service.TransactionService) *topupHandler {
	return &topupHandler{priceService, accountService, transactionService}
}

func (h *topupHandler) handleTopupRequest(c *gin.Context) {
	var request request.CreateTopupRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	price, err := h.priceService.Find()

	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	if request.Harga != price.HargaTopup {
		fmt.Println("Error pricing", err)
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	if ok := helper.ValidateGram(request.Gram); !ok {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	transactionData := model.Transaction{
		ID:           helper.GenShortId(),
		Norek:        request.Norek,
		Gram:         request.Gram,
		HargaTopup:   request.Harga,
		HargaBuyback: price.HargaBuyback,
		Type:         "TOPUP",
		Date:         helper.TimeNow(),
	}

	bytes, err := json.Marshal(transactionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal JSON"})
		return
	}

	producer, err := CreateSyncProducer([]string{"kafka:9092"})
	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}
	defer producer.Close()

	err = SendMessageToKafka(producer, "topup", request.Norek, bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message to Kafka"})
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	mu.Lock()
	responseChannels[request.Norek] = responseCh
	mu.Unlock()

	select {
	case responseMsg := <-responseCh:
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
		delete(responseChannels, request.Norek)
		mu.Unlock()
		error := helper.APIResponseError(true, helper.GenShortId(), "timeout waiting for response")
		c.JSON(http.StatusBadRequest, error)
	}
}
