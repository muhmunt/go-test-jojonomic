package main

import (
	"cek-harga-service/config"
	"cek-harga-service/repository"
	"cek-harga-service/service"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

var (
	responseChannels map[string]chan *sarama.ConsumerMessage
	mu               sync.Mutex
)

func main() {
	db := config.InitDB()

	repository := repository.NewPrice(db)
	service := service.NewService(repository)
	handler := NewPrice(service)

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/check-harga", handler.handleGetPriceRequest)

	if err := router.Run(":8084"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
