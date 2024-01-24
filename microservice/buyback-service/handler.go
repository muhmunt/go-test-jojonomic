package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func handlePingRequest(c *gin.Context) {
	requestID := uuid.New().String()

	message := MyMessage{
		ID:    requestID,
		Name:  "BUYBACK",
		Value: "BUYBACK",
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal JSON"})
		return
	}

	producer, err := CreateSyncProducer([]string{"kafka:9092"})
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		c.JSON(500, gin.H{"error": "failed to create Kafka producer"})
		return
	}
	defer producer.Close()

	err = SendMessageToKafka(producer, "buyback", requestID, bytes)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to send message to Kafka"})
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	mu.Lock()
	responseChannels[requestID] = responseCh
	mu.Unlock()

	select {
	case responseMsg := <-responseCh:
		c.JSON(200, gin.H{"message": string(responseMsg.Value)})
	case <-time.After(10 * time.Second):
		mu.Lock()
		delete(responseChannels, requestID)
		mu.Unlock()
		c.JSON(500, gin.H{"error": "timeout waiting for response"})
	}
}
