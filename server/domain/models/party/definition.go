package party

import (
	proto "soa.mafia-game/proto/mafia-game"
)

const (
	PARTY_SIZE = 5
)

type PartiesDistribution struct {
	party              map[string]int
	non_full_party_ids []int
	party_size         []int
	party_set          int
	roles              map[string]proto.Roles
}
