package event

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var TopicName string = "icEvent"

type KafkaEvent struct {
	producer     *kafka.Producer
	deliveryChan chan kafka.Event
	topic        string
}

func NewKafkaEvent() *KafkaEvent {
	return &KafkaEvent{}
}

func (k *KafkaEvent) Driver() string {
	return "kafka"
}
func (k *KafkaEvent) MaxId() uint64 {
	return uint64(0)
}
func (k *KafkaEvent) Add(e *Event) error {
	data, _ := json.Marshal(e)
	var msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &TopicName,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}
	k.producer.Produce(msg, k.deliveryChan)
	return nil
}
func (k *KafkaEvent) Fetch(eventId uint64) *Event {
	return nil
}
func (k *KafkaEvent) FetchGte(eventId uint64, limit int32) []*Event {
	return nil
}
