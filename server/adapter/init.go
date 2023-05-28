package mafia_server

import (
	"sync"

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

	callbacks_guard sync.Mutex
	user_callbacks  map[string]func()
	brokerServers   string
}

func New(brokerServers string) (*ServerAdapter, error) {
	return &ServerAdapter{
		game:        mafia_domain.NewGame(),
		connections: make(map[string]chan mafia_domain.Event),

		victims:       make(map[string]chan string),
		moved_players: make(map[int]int),

		user_callbacks: make(map[string]func()),

		brokerServers: brokerServers,
	}, nil
}
