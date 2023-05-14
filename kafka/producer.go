package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderPlacer struct {
	producer   *kafka.Producer
	topic      string
	deliverych chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:   p,
		topic:      topic,
		deliverych: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) PlaceOrder(orderType string, size int) error {

	var (
		format  = fmt.Sprintf("%s - %d", orderType, size)
		payload = []byte(format)
	)

	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &op.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	},
		op.deliverych,
	)

	if err != nil {
		log.Fatal(err)
	}

	//wait for the delivery report

	event := <-op.deliverych
	message := event.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		return message.TopicPartition.Error
	}
	fmt.Printf("placed order on the queue %s\n", format)
	return nil

}
