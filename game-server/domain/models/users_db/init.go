package usersdb

import (
	"sync"

	"soa.mafia-game/game-server/domain/models/user"
)

type UsersStorage struct {
	users map[string]user.User
	guard sync.Mutex
}

func New() *UsersStorage {
	return &UsersStorage{
		users: make(map[string]user.User),
	}
}
