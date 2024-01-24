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

	responseChannels = make(map[string]chan *sarama.ConsumerMessage)

	producer, err := CreateSyncProducer([]string{"kafka:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := CreateConsumer([]string{"kafka:9092"})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := CreatePartitionConsumer(consumer, "pong")
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	go ConsumeMessages(partConsumer)

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/check-harga", handler.handleGetPriceRequest)

	if err := router.Run(":8084"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
