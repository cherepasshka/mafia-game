package party

import (
	proto "soa.mafia-game/proto/mafia-game"
)

func New() PartiesDistribution {
	distribution := PartiesDistribution{
		party_size:    make(map[int]int),
		current_party: 0,
		party:         make(map[string]int),
		roles:         make(map[string]proto.Roles),
	}
	return distribution
}
