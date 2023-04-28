package mafia_domain

import (
	"soa.mafia-game/server/domain/models/party"
)

func NewGame() *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        make([]string, 0),
		is_alive:     make(map[string]bool),

		votes_cnt:    make(map[int]map[string]int),
		voted:        make(map[int]int32),
		ghost:        make(map[string]chan string),
		RecentVictim: make(map[int]string),
	}
	return game
}
