package usersdb

import (
	"soa.mafia-game/game-server/domain/models/user"
)

func (storage *UsersStorage) Set(key string, user user.User) {
	storage.guard.Lock()
	defer storage.guard.Unlock()
	storage.users[key] = user
}

func (storage *UsersStorage) Delete(key string) {
	storage.guard.Lock()
	defer storage.guard.Unlock()
	delete(storage.users, key)
}

func (storage *UsersStorage) Get(key string) (user.User, bool) {
	storage.guard.Lock()
	defer storage.guard.Unlock()
	user, exists := storage.users[key]
	return user, exists
}
