package chat

import (
	"context"
	"fmt"
	"log"
	"strings"
	// "time"

	"github.com/segmentio/kafka-go"
	"soa.mafia-game/client/internal/utils/console"
	kafka_service "soa.mafia-game/kafka-help"
)

func (service *ChatService) Start(user_login, sessionId string, partition int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.Listen(ctx, sessionId, partition)

	fmt.Printf("To stop messaging type `exit`")
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
		// MinBytes:  10e3, // 10KB
		// MaxBytes:  10e6, // 10MB
		// MaxWait:   time.Millisecond * 10,
	})
	reader.SetOffset(0)

	for {
		
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		fmt.Printf("%v: %v\n", message.Key, message.Value)
		select {
			case <-ctx.Done():
				return 
		}
	}
}