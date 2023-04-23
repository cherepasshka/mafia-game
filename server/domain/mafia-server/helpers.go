package mafia_domain

import (
	// "fmt"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
)

func (g *MafiaGame) AddPlayer(login string) (bool, Event) {
	/*
		check if user login is unique for party
		add player to party
	*/

	for _, user_login := range g.users {
		if user_login == login {
			return false, Event{}
		}
	}
	g.users = append(g.users, login)
	g.distribution.AddPlayer(login)

	g.Events = append(g.Events, Event{User: login, Status: proto.State_connected, Time: time.Now()})
	return true, g.Events[len(g.Events)-1]
}

func (g *MafiaGame) RemovePlayer(login string) (bool, Event) {
	/*
		check if user login is unique for party
		add player to party
	*/
	user_id := -1
	for i, user_login := range g.users {
		if user_login == login {
			user_id = i
		}
	}
	if user_id == -1 {
		return false, Event{}
	}
	g.users[user_id] = g.users[len(g.users)-1]
	g.users = g.users[:len(g.users)-1] // delete user

	g.distribution.RemovePlayer(login)

	g.Events = append(g.Events, Event{User: login, Status: proto.State_left, Time: time.Now()})
	return true, g.Events[len(g.Events)-1]
}
