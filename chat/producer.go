package chat

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO
func GetNewProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:9092",
		//"client.id":         "prod-1TODO",
		"acks":              "all",
	})
	return producer, err
}

func Produce(key string, value string, topic string, producer *kafka.Producer) {
	deliveryChan := make(chan kafka.Event, 1)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: []byte(value),
		
	}, deliveryChan)

	if err != nil {
		panic(err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("delivery failed %v \n", m.TopicPartition.Error)
	} else {
		fmt.Printf("message delivered topic: %s | key: %s\n", topic, string(key))
	}

	close(deliveryChan)
}

func CreateTopic(admin *kafka.AdminClient, topicName string, numPartitions int) error {
	topic := kafka.TopicSpecification{
		Topic: topicName,
		NumPartitions: numPartitions,
	}
	_, err := admin.CreateTopics(context.Background(), []kafka.TopicSpecification{topic,})
	return err
}

func DeleteTopic(admin *kafka.AdminClient, topicName string) error {
	_, err := admin.DeleteTopics(context.Background(), []string{topicName,})
	return err
}