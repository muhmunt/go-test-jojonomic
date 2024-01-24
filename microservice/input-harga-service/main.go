package main

import (
	"input-harga-service/config"
	"input-harga-service/repository"
	"input-harga-service/service"
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

	partConsumer, err := CreatePartitionConsumer(consumer, "input-result")
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	go ConsumeMessages(partConsumer)

	router := gin.Default()

	api := router.Group("/api")

	api.POST("/input-harga", handler.handleInputPriceRequest)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
