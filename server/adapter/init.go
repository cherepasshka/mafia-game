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
	mut         sync.Mutex

	// TODO
	start_next_day map[string]chan bool
	cnt            int
}

func New() *ServerAdapter {
	return &ServerAdapter{
		game:        mafia_domain.NewGame(),
		connections: make(map[string]chan mafia_domain.Event),

		start_next_day: make(map[string]chan bool),
		cnt:            0,
	}
}
