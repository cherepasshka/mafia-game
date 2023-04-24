package mafia_server

import (
	// "sync"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-server"
)

type mafiaServer struct {
	proto.UnimplementedMafiaServiceServer
	game     *mafia_domain.MafiaGame
	channels map[string]chan mafia_domain.Event

	// TODO
	ready map[string]chan bool
	cnt   int

	// wg sync.WaitGroup

}

func New() *mafiaServer {
	return &mafiaServer{
		game:     mafia_domain.NewGame(),
		channels: make(map[string]chan mafia_domain.Event),

		ready: make(map[string]chan bool),
		cnt:   0,
	}
}
