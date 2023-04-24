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

	Events []Event
}

func NewGame() *MafiaGame {
	game := &MafiaGame{
		distribution: party.New(),
		users:        make([]string, 0),
	}
	return game
}
