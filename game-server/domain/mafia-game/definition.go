package mafia_domain

import (
	"sync"
	"time"

	"soa.mafia-game/game-server/domain/models/party"
	usersdb "soa.mafia-game/game-server/domain/models/users_db"
	proto "soa.mafia-game/proto/mafia-game"
)

type Event struct {
	User             string
	Status           proto.State
	Time             time.Time
	SessionReadiness bool
}

type MafiaGame struct {
	distribution party.PartiesDistribution
	users        *usersdb.UsersStorage
	is_alive     map[string]bool

	guard        sync.Mutex
	ghost        map[string]chan string
	votes_cnt    map[int]map[string]int
	voted        map[int]int32
	RecentVictim map[int]string
	Events       []Event
}
