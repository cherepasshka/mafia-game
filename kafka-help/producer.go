package kafka_service

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO
func GetNewProducer(brokerServers string) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerServers, //"kafka1:9092",
		"client.id":         "clients",
		"acks":              "all",
	})
	log.Printf("After creating producer with brokers %v\n", brokerServers)
	return producer, err
}

func Produce(key string, value string, topic string, partition int32, producer *kafka.Producer) error {
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
		},
		Key:   []byte(key),
		Value: []byte(value),
	}, deliveryChan)

	if err != nil {
		log.Printf("FAILED TO PRODUCE %v\n", err)
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Printf("FAILED TO PRODUCE %v\n", err)
		return m.TopicPartition.Error
	} else {
		fmt.Printf("message `%s` delivered topic: %s | key: %s| part %v\n", value, topic, string(key), partition)
	}

	log.Printf("ALL FINE -___-")
	return nil
}

func CreateTopic(admin *kafka.AdminClient, topicName string, numPartitions int) error {
	log.Printf("CREATE TOPIC %v", topicName)
	topic := kafka.TopicSpecification{
		Topic:         topicName,
		NumPartitions: numPartitions,
	}
	_, err := admin.CreateTopics(context.Background(), []kafka.TopicSpecification{topic})
	return err
}

func DeleteTopic(admin *kafka.AdminClient, topicName string) error {
	_, err := admin.DeleteTopics(context.Background(), []string{topicName})
	return err
}
