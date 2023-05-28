package party

import (
	"math/rand"
	"sort"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
)

func (d *PartiesDistribution) AddPlayer(user_login string) {
	d.party_mutex.Lock()
	defer d.party_mutex.Unlock()
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
	d.party_mutex.Lock()
	defer d.party_mutex.Unlock()
	party, ok := d.party[user_login]
	if !ok {
		return
	}
	if d.party_size[party] == PARTY_SIZE {
		d.non_full_party_ids = append(d.non_full_party_ids, party)
	}
	d.party_size[party]--
	delete(d.party, user_login)
}

func (d *PartiesDistribution) GetUserParty(user_login string) int {
	d.party_mutex.Lock()
	defer d.party_mutex.Unlock()
	return d.party[user_login]
}

func (d *PartiesDistribution) GetPartySize(party int) int {
	d.party_mutex.Lock()
	defer d.party_mutex.Unlock()
	return d.party_size[party]
}

func (d *PartiesDistribution) IsFull(party int) bool {
	return d.GetPartySize(party) == PARTY_SIZE
}

func (d *PartiesDistribution) GetParty(party int) []string {
	d.party_mutex.Lock()
	defer d.party_mutex.Unlock()
	result := make([]string, PARTY_SIZE)
	ind := 0
	for user := range d.party {
		if d.party[user] == party {
			result[ind] = user
			ind++
		}
	}
	result = result[:ind]
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

	d.roles_mutex.Lock()
	defer d.roles_mutex.Unlock()
	for i := 0; i < CIVILIANS; i++ {
		d.roles[members[i]] = proto.Roles_Civilian
	}
	for i := CIVILIANS; i < CIVILIANS+MAFIAS; i++ {
		d.roles[members[i]] = proto.Roles_Mafia
	}
	for i := CIVILIANS + MAFIAS; i < CIVILIANS+MAFIAS+COMMISSIONERS; i++ {
		d.roles[members[i]] = proto.Roles_Commissioner
	}
	return true
}

func (d *PartiesDistribution) GetRole(user string) proto.Roles {
	d.roles_mutex.Lock()
	defer d.roles_mutex.Unlock()
	return d.roles[user]
}
