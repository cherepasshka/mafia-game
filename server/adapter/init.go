package mafia_server

import (
	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-server"
)

type mafiaServer struct {
	proto.UnimplementedMafiaServiceServer
	game     *mafia_domain.MafiaGame
	channels []chan mafia_domain.Event
}

func New() *mafiaServer {
	return &mafiaServer{
		game:     mafia_domain.NewGame(),
		channels: make([]chan mafia_domain.Event, 0),
	}
}
