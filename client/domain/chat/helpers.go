package chat

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"

	"soa.mafia-game/client/internal/utils/console"
	kafka_service "soa.mafia-game/kafka-help"
)

func (service *ChatService) Start(user_login, sessionId string, partition int32) {
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
		if msg == "exit" {
			return
		}
		kafka_service.Produce(user_login, msg, user_login, 0, service.producer)
	}
}

func (service *ChatService) Listen(ctx context.Context, sessionId string, partition int32) {
	brokers := strings.Split(service.brokerServers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     sessionId,
		Partition: int(partition),
	})
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
