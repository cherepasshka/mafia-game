package party

import (
	proto "soa.mafia-game/proto/mafia-game"
)

func New() PartiesDistribution {
	distribution := PartiesDistribution{
		party_size:         make(map[int]int),
		current_party: 0,
		// non_full_party_ids: make([]int, 1),
		// party_set:          1,
		party:              make(map[string]int),
		roles:              make(map[string]proto.Roles),
	}
	// distribution.non_full_party_ids[0] = 0
	// distribution.party_size[0] = 0
	return distribution
}
