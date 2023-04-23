package party

const (
	PARTY_SIZE = 4
)

type PartiesDistribution struct {
	party              map[string]int
	non_full_party_ids []int
	party_size         []int
	// users              []string
	party_set int
}

func New() PartiesDistribution {
	distribution := PartiesDistribution{
		party_size:         make([]int, 1),
		non_full_party_ids: make([]int, 1),
		party_set:          1,
		party:              make(map[string]int),
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

func (d *PartiesDistribution) GetParty(user_login string) int {
	return d.party[user_login]
}
