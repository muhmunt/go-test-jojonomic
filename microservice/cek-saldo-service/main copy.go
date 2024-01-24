package main

import (
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
	router.GET("/ping", handlePingRequest)

	if err := router.Run(":8086"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
