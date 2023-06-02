package usersdb

import (
	"sync"

	"soa.mafia-game/game-server/domain/models/party"
	"soa.mafia-game/game-server/domain/models/user"
)

type Storage struct {
	party.PartiesDistribution

	users       map[string]user.User
	users_guard sync.Mutex

	// distribution party.PartiesDistribution
	// distrib_guard sync.Mutex
}

func New() *Storage {
	return &Storage{
		users:               make(map[string]user.User),
		PartiesDistribution: party.New(),
	}
}
