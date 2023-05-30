package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"

	// "soa.mafia-game/client/internal/utils/console"
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

	go service.Listen(ctx, user_login, sessionId, partition)

	color.Black("To stop messaging type `exit`")
	for {
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := line[:len(line)-1]
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

func (service *ChatService) Listen(ctx context.Context, user_login, sessionId string, partition int32) {
	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error while creating chat: %v", err)
	}
	defer func() {
		admin.Close()
	}()
	_ = admin.CreateTopic(sessionId, &sarama.TopicDetail{NumPartitions: 4, ReplicationFactor: 1}, false)

	brokers := strings.Split(service.brokerServers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     sessionId,
		Partition: int(partition),
	})
	defer reader.Close()
	reader.SetOffset(0)

	number := make(map[string]int)
	ind := 0
	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(context.Canceled, err) {
				return
			} else {
				log.Printf("%v\n", err)
			}
		}
		user := string(message.Key)
		if _, exists := number[user]; !exists {
			number[user] = ind
			ind++
		}
		if user != user_login {
			colorfulPrint(fmt.Sprintf("%v: %v", user, string(message.Value)), number[user])
			// color.Black("%v says %v", , string(message.Value))
		}
	}
}

func colorfulPrint(value string, number int) {
	if number == 0 {
		color.Blue(value)
	} else if number == 1 {
		color.Green(value)
	} else if number == 2 {
		color.Red(value)
	} else {
		color.Magenta(value)
	}
}
