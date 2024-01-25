package main

import (
	"buyback-service/config"
	"buyback-service/repository"
	"buyback-service/service"
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

	priceRepository := repository.NewPrice(db)
	priceService := service.NewPrice(priceRepository)
	accountRepository := repository.NewAccount(db)
	accountService := service.NewAccount(accountRepository)
	transactionRepository := repository.NewTransaction(db)
	transactionService := service.NewTransaction(transactionRepository)
	buybackHandler := NewBuyback(priceService, accountService, transactionService)

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

	partConsumer, err := CreatePartitionConsumer(consumer, "buyback-result")
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	go ConsumeMessages(partConsumer)

	router := gin.Default()
	api := router.Group("api")
	api.POST("/buyback", buybackHandler.handleBuybackRequest)

	if err := router.Run(":8087"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
