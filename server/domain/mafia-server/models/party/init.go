package party

import (
	"math/rand"
	proto "soa.mafia-game/proto/mafia-game"
	"sort"
	"time"
)

const (
	PARTY_SIZE = 4
)

type PartiesDistribution struct {
	party              map[string]int
	non_full_party_ids []int
	party_size         []int
	// users              []string
	party_set int
	roles     map[string]proto.Roles
}

func New() PartiesDistribution {
	distribution := PartiesDistribution{
		party_size:         make([]int, 1),
		non_full_party_ids: make([]int, 1),
		party_set:          1,
		party:              make(map[string]int),
		roles:              make(map[string]proto.Roles),
	}
	distribution.non_full_party_ids[0] = 0
	distribution.party_size[0] = 0
	return distribution
}

func (d *PartiesDistribution) AddPlayer(user_login string) {
	id := len(d.non_full_party_ids) - 1
	party := d.non_full_party_ids[id]
	d.party[user_login] = party
	d.party_size[party]++
	if d.party_size[party] == PARTY_SIZE {
		d.non_full_party_ids[id] = d.party_set
		d.party_set++
	}
}

func (d *PartiesDistribution) RemovePlayer(user_login string) {
	party := d.party[user_login]
	if d.party_size[party] == PARTY_SIZE {
		d.non_full_party_ids = append(d.non_full_party_ids, party)
	}
	d.party_size[party]--
	delete(d.party, user_login)
}

func (d *PartiesDistribution) GetUserParty(user_login string) int {
	return d.party[user_login]
}

func (d *PartiesDistribution) GetPartySize(party int) int {
	return d.party_size[party]
}

func (d *PartiesDistribution) IsFull(party int) bool {
	return d.GetPartySize(party) == PARTY_SIZE
}

func (d *PartiesDistribution) GetParty(party int) []string {
	// could be smarter but im lazy
	result := make([]string, PARTY_SIZE)
	ind := 0
	for user := range d.party {
		if d.party[user] == party {
			result[ind] = user
			ind++
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func (d *PartiesDistribution) DistributeRoles(party int) bool {

	members := d.GetParty(party)
	if len(members) != PARTY_SIZE {
		return false
	}
	sort.Slice(members, func(i, j int) bool { return members[i] < members[j] })
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(members), func(i, j int) { members[i], members[j] = members[j], members[i] })

	d.roles[members[0]] = proto.Roles_Civilian
	d.roles[members[1]] = proto.Roles_Civilian
	d.roles[members[2]] = proto.Roles_Mafia
	d.roles[members[3]] = proto.Roles_Commissioner
	return true
}

func (d *PartiesDistribution) GetRole(user string) proto.Roles {
	return d.roles[user]
}
