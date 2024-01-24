package main

import (
	"log"

	"github.com/IBM/sarama"
)

func CreateSyncProducer(brokers []string) (sarama.SyncProducer, error) {

	producer, err := sarama.NewSyncProducer(brokers, nil)
	return producer, err
}

func CreateConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	return consumer, err
}

func CreatePartitionConsumer(consumer sarama.Consumer, topic string) (sarama.PartitionConsumer, error) {
	partConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	return partConsumer, err
}

func ConsumeMessages(partConsumer sarama.PartitionConsumer) {
	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}
			responseID := string(msg.Key)
			handleResponse(responseID, msg)
		}
	}
}

func SendMessageToKafka(producer sarama.SyncProducer, topic, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := producer.SendMessage(msg)
	return err
}

func handleResponse(responseID string, msg *sarama.ConsumerMessage) {
	mu.Lock()
	defer mu.Unlock()

	ch, exists := responseChannels[responseID]
	if exists {
		ch <- msg
		delete(responseChannels, responseID)
	}
}
