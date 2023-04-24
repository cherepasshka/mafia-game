package mafia_domain

import (
	"soa.mafia-game/server/domain/models/party"
)

func NewGame() *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        make([]string, 0),
		is_alive:     make(map[string]bool),
	}
	return game
}
