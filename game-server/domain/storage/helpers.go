package usersdb

import (
	"soa.mafia-game/game-server/domain/models/user"
)

func (storage *Storage) SetUser(key string, user user.User) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	storage.users[key] = user
}

func (storage *Storage) DeleteUser(key string) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	delete(storage.users, key)
}

func (storage *Storage) GetUser(key string) (user.User, bool) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	user, exists := storage.users[key]
	return user, exists
}
