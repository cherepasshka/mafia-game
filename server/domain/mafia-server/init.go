package mafia_domain

import (
	// domain "soa.mafia-game/server/domain"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
	"soa.mafia-game/server/domain/mafia-server/models/party"
)

type Event struct {
	User             string
	Status           proto.State
	Time             time.Time
	SessionReadiness bool
}

type MafiaGame struct {
	distribution party.PartiesDistribution
	users        []string
	is_alive     map[string]bool
	RecentVictim string
	Events       []Event
}

func NewGame() *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        make([]string, 0),
		is_alive:     make(map[string]bool),
	}
	return game
}
