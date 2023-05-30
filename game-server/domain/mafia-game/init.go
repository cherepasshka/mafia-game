package mafia_domain

import (
	"soa.mafia-game/game-server/domain/models/party"
	"soa.mafia-game/game-server/domain/models/user"
)

func NewGame() *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        make(map[string]user.User),
		is_alive:     make(map[string]bool),

		votes_cnt:    make(map[int]map[string]int),
		voted:        make(map[int]int32),
		ghost:        make(map[string]chan string),
		RecentVictim: make(map[int]string),
	}
	return game
}
