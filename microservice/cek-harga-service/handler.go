package main

import (
	"cek-harga-service/helper"
	"cek-harga-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type priceHandler struct {
	priceService service.Service
}

func NewPrice(service service.Service) *priceHandler {
	return &priceHandler{service}
}

func (h *priceHandler) handleGetPriceRequest(c *gin.Context) {
	price, err := h.priceService.Find()

	if err != nil {
		errors := helper.ValidationFormatError(err)
		error := helper.APIResponseError(true, helper.GenShortId(), errors)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"reff_id": helper.GenShortId(),
		"data":    price,
	})
}
