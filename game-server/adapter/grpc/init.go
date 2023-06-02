package mafia_server

import (
	"sync"

	mafia_domain "soa.mafia-game/game-server/domain/mafia-game"
	usersdb "soa.mafia-game/game-server/domain/models/storage"
	proto "soa.mafia-game/proto/mafia-game"
)

type ServerAdapter struct {
	proto.UnimplementedMafiaServiceServer
	game        *mafia_domain.MafiaGame
	connections map[string]chan mafia_domain.Event
	guard       sync.Mutex
	conn_guard  sync.Mutex

	victims       map[string]chan string
	moved_players map[int]int
}

func New(users *usersdb.Storage, brokerServers string) (*ServerAdapter, error) {
	return &ServerAdapter{
		game:        mafia_domain.NewGame(users),
		connections: make(map[string]chan mafia_domain.Event),

		victims:       make(map[string]chan string),
		moved_players: make(map[int]int),
	}, nil
}
