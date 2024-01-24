package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"topup-storage-service/config"
	"topup-storage-service/model"
	"topup-storage-service/repository"
	"topup-storage-service/service"

	"github.com/IBM/sarama"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	db := config.InitDB()

	priceRepository := repository.NewPrice(db)
	priceService := service.NewPrice(priceRepository)
	accountRepository := repository.NewAccount(db)
	accountService := service.NewAccount(accountRepository)
	transactionRepository := repository.NewTransaction(db)
	transactionService := service.NewTransaction(transactionRepository)

	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition("topup", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting")
				return
			}

			var transaction model.Transaction
			err := json.Unmarshal(msg.Value, &transaction)

			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			// getSaldo := accountService.
			// transaction.SaldoTerakhir =
			_, err = transactionService.StoreTransaction(transaction)

			resp := &sarama.ProducerMessage{
				Topic: "top",
				Key:   sarama.StringEncoder(transaction.ID),
				Value: sarama.StringEncoder("responseText"),
			}

			_, _, err = producer.SendMessage(resp)
			if err != nil {
				log.Printf("Failed to send message to Kafka: %v", err)
			}
		case <-signals:
			log.Println("Received interrupt signal, exiting")
		}
	}
}
