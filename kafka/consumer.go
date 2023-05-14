package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeMessages(consumer *kafka.Consumer, broadcast chan<- []byte) {
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error while consuming message: %s\n", err)
			continue
		}
		broadcast <- message.Value
	}
}
