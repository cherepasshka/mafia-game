package game

import (
	"soa.mafia-game/client/domain/models/user"
)

type Game struct {
	player    models.User
	players   []string
	alive     []string
	// sessionId int32
	// partition int32
}
