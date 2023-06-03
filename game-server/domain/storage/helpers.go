package usersdb

import (
	"time"

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

func (storage *Storage) SetEnterTime(key string, enterTime time.Time) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	user := storage.users[key]
	user.ActiveSession = true
	user.LastSessionEnter = enterTime

	storage.users[key] = user
}

func (storage *Storage) LeaveGameSession(key string) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	user := storage.users[key]
	user.ActiveSession = false
	user.TotalGameTime += time.Now().Sub(user.LastSessionEnter)

	storage.users[key] = user
}

func (storage *Storage) IncGameCnt(key string, isVictory bool) {
	storage.users_guard.Lock()
	defer storage.users_guard.Unlock()
	user := storage.users[key]
	user.SessionsCnt++
	if isVictory {
		user.VictoriesCnt++
	}

	storage.users[key] = user
}
