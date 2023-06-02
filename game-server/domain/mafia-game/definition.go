package mafia_domain

import (
	"sync"
	"time"

	usersdb "soa.mafia-game/game-server/domain/models/storage"
	proto "soa.mafia-game/proto/mafia-game"
)

type Event struct {
	User             string
	Status           proto.State
	Time             time.Time
	SessionReadiness bool
}

type MafiaGame struct {
	storage  *usersdb.Storage
	is_alive map[string]bool

	guard        sync.Mutex
	ghost        map[string]chan string
	votes_cnt    map[int]map[string]int
	voted        map[int]int32
	RecentVictim map[int]string
	Events       []Event
}
