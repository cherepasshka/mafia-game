package mafia_domain

import (
	"time"

	party_model "soa.mafia-game/game-server/domain/models/party"
	"soa.mafia-game/game-server/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

func (game *MafiaGame) AddPlayer(login string) (bool, Event) {
	game.guard.Lock()
	defer game.guard.Unlock()
	_, exists := game.storage.GetUser(login)
	if exists {
		return false, Event{}

	}
	new_user := user.New(login)
	game.storage.SetUser(login, new_user)

	game.Events = append(game.Events, Event{User: login, Status: proto.State_connected, Time: time.Now()})
	return true, game.Events[len(game.Events)-1]
}

func (game *MafiaGame) RemovePlayer(login string) (bool, Event) {
	game.guard.Lock()
	defer game.guard.Unlock()
	_, exists := game.storage.GetUser(login)
	if !exists {
		return false, Event{}
	}
	game.storage.DeleteUser(login)

	game.storage.RemovePlayer(login)

	game.Events = append(game.Events, Event{User: login, Status: proto.State_left, Time: time.Now()})
	return true, game.Events[len(game.Events)-1]
}

func (game *MafiaGame) SessionReady(user string) bool {
	return game.storage.IsFull(game.storage.GetUserParty(user))
}

func (game *MafiaGame) GetParty(user string) int {
	return game.storage.GetUserParty(user)
}

func (game *MafiaGame) GetPartition(user string) int {
	members := game.storage.GetParty(game.GetParty(user))
	for i := range members {
		if members[i] == user {
			return i
		}
	}
	return -1
}

func (game *MafiaGame) GetRole(user string) proto.Roles {
	return game.storage.GetRole(user)
}

func (game *MafiaGame) DistributeRoles(party int) bool {
	members := game.GetMembers(party)
	game.guard.Lock()
	defer game.guard.Unlock()
	for _, member := range members {
		game.is_alive[member] = true
	}
	return game.storage.DistributeRoles(party)
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
	return game.storage.GetParty(party)
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

func (game *MafiaGame) ExitGame(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	isVictory := false
	winner := game.Winner(game.GetParty(user_login))
	if game.GetRole(user_login) == proto.Roles_Mafia {
		isVictory = (winner == proto.Roles_Mafia)
	} else {
		isVictory = (winner == proto.Roles_Civilian)
	}

	game.storage.IncGameCnt(user_login, isVictory)
	game.storage.LeaveGameSession(user_login)
}

func (game *MafiaGame) ExitSession(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.storage.AddPlayer(user_login)
}

func (game *MafiaGame) EnterSession(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.storage.AddPlayer(user_login)
}

func (game *MafiaGame) EnterGame(user_login string) {
	game.guard.Lock()
	defer game.guard.Unlock()
	game.storage.SetEnterTime(user_login, time.Now())
}
