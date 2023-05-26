package mafia_server

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"soa.mafia-game/chat"
	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-game"
)

type ServerAdapter struct {
	proto.UnimplementedMafiaServiceServer
	game        *mafia_domain.MafiaGame
	connections map[string]chan mafia_domain.Event
	guard       sync.Mutex
	conn_guard  sync.Mutex

	victims       map[string]chan string
	moved_players map[int]int

	producer *kafka.Producer
}

func New() *ServerAdapter {
	producer, _ := chat.GetNewProducer() // TODO
	
	return &ServerAdapter{
		game:        mafia_domain.NewGame(),
		connections: make(map[string]chan mafia_domain.Event),

		victims:       make(map[string]chan string),
		moved_players: make(map[int]int),

		producer: producer,
	}
}
