package chat

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"

	"soa.mafia-game/client/internal/utils/console"
	kafka_service "soa.mafia-game/kafka-help"
)

func (service *ChatService) Start(user_login, sessionId string, partition int32, isGhost bool) {
	producer, _ := kafka_service.GetNewProducer("localhost:9092")
	defer producer.Close()
	if isGhost {
		kafka_service.Produce(sessionId, "exit", user_login, 0, producer)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go service.Listen(ctx, sessionId, partition)

	color.Black("To stop messaging type `exit`")
	for {
		msg, err := console.Ask(">")
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		kafka_service.Produce(sessionId, msg, user_login, 0, producer)
		if msg == "exit" {
			return
		}
	}
}

func (service *ChatService) Listen(ctx context.Context, sessionId string, partition int32) {
	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error while creating chat: %v", err)
	}
	defer func() {
		admin.Close()
	}()
	_ = admin.CreateTopic(sessionId, &sarama.TopicDetail{NumPartitions: 4, ReplicationFactor: 1}, false)
	// if err != nil && !errors.Is(sarama.ErrTopicAlreadyExists, err) {
	//     log.Fatalf("Failed to open chat: %v; %v", err, errors.Is(sarama.ErrTopicAlreadyExists, err))
	// }

	brokers := strings.Split(service.brokerServers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     sessionId,
		Partition: int(partition),
	})
	defer reader.Close()
	reader.SetOffset(0)

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(context.Canceled, err) {
				return
			} else {
				log.Printf("%v\n", err)
			}
		}
		color.Black("%v says %v", string(message.Key), string(message.Value))
	}
}
