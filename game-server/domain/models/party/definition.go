package party

import (
	"sync"

	proto "soa.mafia-game/proto/mafia-game"
)

const (
	PARTY_SIZE    = 4
	CIVILIANS     = 2
	MAFIAS        = 1
	COMMISSIONERS = 1
)

type PartiesDistribution struct {
	party         map[string]int // guarded by mutex
	party_size    map[int]int    // guarded by mutex
	current_party int
	roles         map[string]proto.Roles //guarded by roles_mutex
	party_mutex   sync.Mutex
	roles_mutex   sync.Mutex
}
