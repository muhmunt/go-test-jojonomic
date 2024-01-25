package main

import (
	"buyback-service/helper"
	"buyback-service/model"
	"buyback-service/request"
	"buyback-service/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type buybackHandler struct {
	priceService       service.PriceService
	accountService     service.AccountService
	transactionService service.TransactionService
}

func NewBuyback(priceService service.PriceService, accountService service.AccountService, transactionService service.TransactionService) *buybackHandler {
	return &buybackHandler{priceService, accountService, transactionService}
}

func (h *buybackHandler) handleBuybackRequest(c *gin.Context) {
	var request request.CreateTopupRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	price, err := h.priceService.Find()
	account, err := h.accountService.FindById(request.Norek)
	requestGram, err := helper.DecimalFromString(request.Gram)

	if requestGram > account.Saldo {
		fmt.Println("Error gram input", err)
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusBadRequest, error)
		return
	}

	if request.Harga != price.HargaBuyback {
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
		HargaTopup:   price.HargaTopup,
		HargaBuyback: request.Harga,
		Type:         "BUYBACK",
		Date:         helper.TimeNow(),
	}

	bytes, err := json.Marshal(transactionData)
	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	producer, err := CreateSyncProducer([]string{"kafka:9092"})
	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusInternalServerError, error)
		return
	}
	defer producer.Close()

	err = SendMessageToKafka(producer, "buyback", request.Norek, bytes)
	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	mu.Lock()
	responseChannels[request.Norek] = responseCh
	mu.Unlock()

	select {
	case responseMsg := <-responseCh:
		if string(responseMsg.Value) == "false" {
			error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
			c.JSON(http.StatusInternalServerError, error)
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
		error := helper.APIResponseError(true, helper.GenShortId(), "Kafka not ready")
		c.JSON(http.StatusInternalServerError, error)
		return
	}
}
