package mafia_domain

import (
	"time"

	proto "soa.mafia-game/proto/mafia-game"
	party_model "soa.mafia-game/server/domain/models/party"
)

func (game *MafiaGame) AddPlayer(login string) (bool, Event) {
	for _, user_login := range game.users {
		if user_login == login {
			return false, Event{}
		}
	}
	game.guard.Lock()
	defer game.guard.Unlock()
	game.users = append(game.users, login)

	game.Events = append(game.Events, Event{User: login, Status: proto.State_connected, Time: time.Now()})
	return true, game.Events[len(game.Events)-1]
}

func (game *MafiaGame) RemovePlayer(login string) (bool, Event) {
	user_id := -1
	for i, user_login := range game.users {
		if user_login == login {
			user_id = i
		}
	}
	if user_id == -1 {
		return false, Event{}
	}
	game.guard.Lock()
	defer game.guard.Unlock()
	game.users[user_id] = game.users[len(game.users)-1]
	game.users = game.users[:len(game.users)-1] // delete user

	game.distribution.RemovePlayer(login)

	game.Events = append(game.Events, Event{User: login, Status: proto.State_left, Time: time.Now()})
	return true, game.Events[len(game.Events)-1]
}

func (game *MafiaGame) SessionReady(user string) bool {
	return game.distribution.IsFull(game.distribution.GetUserParty(user))
}

func (game *MafiaGame) GetParty(user string) int {
	return game.distribution.GetUserParty(user)
}

func (game *MafiaGame) GetPartition(user string) int {
	members := game.distribution.GetParty(game.GetParty(user))
	for i := range members {
		if members[i] == user {
			return i
		}
	}
	return -1
}

func (game *MafiaGame) GetRole(user string) proto.Roles {
	return game.distribution.GetRole(user)
}

func (game *MafiaGame) DistributeRoles(party int) bool {
	members := game.GetMembers(party)
	game.guard.Lock()
	defer game.guard.Unlock()
	for _, member := range members {
		game.is_alive[member] = true
	}
	return game.distribution.DistributeRoles(party)
}

func (game *MafiaGame) IsPlayerAlive(login string) bool {
	return game.is_alive[login]
}

func (game *MafiaGame) Kill(login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.is_alive[login] = false
}

func (game *MafiaGame) GetAliveMembers(party int) []string {
	members := game.GetMembers(party)
	alive := make([]string, 0, len(members))
	for _, member := range members {
		if game.is_alive[member] {
			alive = append(alive, member)
		}
	}
	return alive
}

func (game *MafiaGame) GetMembers(party int) []string {
	return game.distribution.GetParty(party)
}

func (game *MafiaGame) CountRole(party int, role proto.Roles) int {
	members := game.GetAliveMembers(party)
	cnt := 0
	for _, member := range members {
		if game.GetRole(member) == role {
			cnt++
		}
	}
	return cnt
}

func (game *MafiaGame) IsActive(party int) bool {
	mafia_cnt := game.CountRole(party, proto.Roles_Mafia)
	civilian_cnt := game.CountRole(party, proto.Roles_Civilian)
	return !(mafia_cnt == 0 || civilian_cnt <= mafia_cnt)
}

func (game *MafiaGame) Winner(party int) proto.Roles {
	mafia_cnt := game.CountRole(party, proto.Roles_Mafia)
	civilian_cnt := game.CountRole(party, proto.Roles_Civilian)
	if mafia_cnt == 0 {
		return proto.Roles_Civilian
	}
	if civilian_cnt <= mafia_cnt {
		return proto.Roles_Mafia
	}
	return proto.Roles_Undefined
}

func (game *MafiaGame) VoteFor(voter_login string, guess string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.ghost[voter_login] = make(chan string, 1)
	party := game.GetParty(voter_login)
	cnt, exist := game.votes_cnt[party]
	if !exist {
		game.votes_cnt[party] = make(map[string]int)
		cnt = game.votes_cnt[party]
	}
	game.voted[party] += 1
	if game.IsPlayerAlive(voter_login) {
		cnt[guess] += 1
	}
}

func (game *MafiaGame) WaitForEverybody(user_login string) string {
	party := game.GetParty(user_login)
	game.guard.Lock()
	if game.voted[party] == party_model.PARTY_SIZE {
		game.voted[party] = 0
		ghost := user_login
		for player := range game.votes_cnt[party] {
			if game.votes_cnt[party][player] > game.votes_cnt[party][ghost] {
				ghost = player
			}
		}
		for key := range game.votes_cnt[party] {
			delete(game.votes_cnt[party], key)
		}
		game.is_alive[ghost] = false
		for _, user := range game.GetMembers(party) {
			game.ghost[user] <- ghost
		}
	}
	game.guard.Unlock()
	ghost := <-game.ghost[user_login]
	return ghost
}

func (game *MafiaGame) ExitSession(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.distribution.RemovePlayer(user_login)
}

func (game *MafiaGame) EnterSession(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.distribution.AddPlayer(user_login)
}
