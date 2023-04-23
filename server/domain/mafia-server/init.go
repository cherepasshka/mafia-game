package mafia_domain

import (
	// domain "soa.mafia-game/server/domain"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
)

const (
	party_size = 4
)

type Event struct {
	User   string
	Status proto.State
	Time   time.Time
}

type MafiaGame struct {
	party              map[string]int
	non_full_party_ids []int
	party_size         []int
	users              []string
	party_cnt          int

	Events []Event
}

func NewGame() *MafiaGame {
	game := &MafiaGame{
		party_size:         make([]int, 1),
		non_full_party_ids: make([]int, 1),
		party_cnt:          1,
		party:              make(map[string]int),
	}
	game.non_full_party_ids[0] = 0
	game.party_size[0] = 0
	return game
}
