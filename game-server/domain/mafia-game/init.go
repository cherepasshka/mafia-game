package mafia_domain

import (
	"soa.mafia-game/game-server/domain/models/party"
	usersdb "soa.mafia-game/game-server/domain/models/users_db"
)

func NewGame(users *usersdb.UsersStorage) *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        users,
		is_alive:     make(map[string]bool),

		votes_cnt:    make(map[int]map[string]int),
		voted:        make(map[int]int32),
		ghost:        make(map[string]chan string),
		RecentVictim: make(map[int]string),
	}
	return game
}
