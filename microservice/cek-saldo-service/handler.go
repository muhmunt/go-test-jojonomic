package main

import (
	"cek-saldo-service/formatter"
	"cek-saldo-service/helper"
	"cek-saldo-service/request"
	"cek-saldo-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	accountService service.AccountService
}

func NewAccount(service service.AccountService) *accountHandler {
	return &accountHandler{service}
}

func (h *accountHandler) handleGetAccountRequest(c *gin.Context) {
	var request request.GetAccountRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}
	account, err := h.accountService.FindById(request.Norek)

	if err != nil {
		error := helper.APIResponseError(true, helper.GenShortId(), "kafka not ready")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"reff_id": helper.GenShortId(),
		"data":    formatter.AccountFormatter(account),
	})
}
