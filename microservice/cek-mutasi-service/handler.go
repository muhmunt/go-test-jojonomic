package main

import (
	"cek-mutasi-service/helper"
	"cek-mutasi-service/request"
	"cek-mutasi-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService service.TransactionService
}

func NewTransaction(service service.TransactionService) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) handleGetTransactionRequest(c *gin.Context) {
	var request request.GetTransactionRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}
	transaction, err := h.transactionService.FindTransactionByNorek(request)

	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "kafka not ready")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"reff_id": helper.GenShortId(),
		"data":    transaction,
	})
}
