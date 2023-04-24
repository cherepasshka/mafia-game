package game

import (
	"soa.mafia-game/client/domain/models/user"
)

func New(player models.User, players []string) *Game {
	return &Game{
		player:  player,
		players: players,
		alive:   players,
	}
}
