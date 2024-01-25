package main

import (
	"buyback-storage-service/config"
	"buyback-storage-service/model"
	"buyback-storage-service/repository"
	"buyback-storage-service/service"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	db := config.InitDB()

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

	partConsumer, err := consumer.ConsumePartition("buyback", 0, sarama.OffsetNewest)
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

			getSaldo, err := accountService.FindById(transaction.Norek)
			transaction.SaldoTerakhir = getSaldo.Saldo

			var result string
			result = "true"

			_, err = transactionService.StoreTransaction(transaction)

			_, err = accountService.UpdateOrInsertAccount(transaction)

			if err != nil {
				result = "false"
			}

			resp := &sarama.ProducerMessage{
				Topic: "buyback-result",
				Key:   sarama.StringEncoder(transaction.Norek),
				Value: sarama.StringEncoder(result),
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
