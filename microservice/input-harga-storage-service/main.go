package main

import (
	"encoding/json"
	"input-harga-storage-service/config"
	"input-harga-storage-service/model"
	"input-harga-storage-service/repository"
	"input-harga-storage-service/service"
	"log"

	"github.com/IBM/sarama"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	db := config.InitDB()

	repository := repository.NewPrice(db)
	service := service.NewService(repository)

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

	partConsumer, err := consumer.ConsumePartition("input-harga", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting")
				return
			}

			var price model.Price
			err := json.Unmarshal(msg.Value, &price)

			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			_, err = service.StorePrice(price)

			var result string

			result = "true"
			if err != nil {
				result = "false"
			}

			resp := &sarama.ProducerMessage{
				Topic: "input-result",
				Key:   sarama.StringEncoder(price.AdminID),
				Value: sarama.StringEncoder(result),
			}

			_, _, err = producer.SendMessage(resp)
			if err != nil {
				log.Printf("Failed to send message to Kafka: %v", err)
			}
		}
	}
}
