package mafia_domain

import (
	"sync"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
	"soa.mafia-game/server/domain/models/party"
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

	mut          sync.Mutex
	ghost        map[string]chan string
	votes_cnt    map[int]map[string]int
	voted        map[int]int32
	RecentVictim string // TODO: for each party
	Events       []Event
}
