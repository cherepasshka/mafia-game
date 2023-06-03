package mafia_domain

import (
	usersdb "soa.mafia-game/game-server/domain/storage"
)

func NewGame(storage *usersdb.Storage) *MafiaGame {
	game := &MafiaGame{
		storage:  storage,
		is_alive: make(map[string]bool),

		votes_cnt:    make(map[int]map[string]int),
		voted:        make(map[int]int32),
		ghost:        make(map[string]chan string),
		RecentVictim: make(map[int]string),
	}
	return game
}
