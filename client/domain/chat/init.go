package chat

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"soa.mafia-game/kafka-help"
)

type ChatService struct {
	producer      *kafka.Producer // to remove
	admin         *kafka.AdminClient // to remove
	brokerServers string
}

func New(brokerServers string) (*ChatService, error) {
	log.Printf("in chat.New\n")
	producer, err := kafka_service.GetNewProducer(brokerServers)
	if err != nil {
		return nil, err
	}
	admin, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		return nil, err
	}
	service := &ChatService{
		producer:      producer,
		admin:         admin,
		brokerServers: brokerServers, // HERE SINGLE
	}
	return service, nil
}
