package models

import (
	"fmt"

	proto "soa.mafia-game/proto/mafia-game"
)

type LiveStatus int

const (
	Alive LiveStatus = 0
	Dead  LiveStatus = 1
)

type User interface {
	GetLogin() string
	GetRole() proto.Roles
	GetStatus() LiveStatus
	SetStatus(LiveStatus)
	MakeNightMove(proto.MafiaServiceClient) error
}

type BaseUser struct {
	Login  string
	Status LiveStatus
}

func (user *BaseUser) GetLogin() string {
	return user.Login
}

func (user *BaseUser) GetRole() proto.Roles {
	return proto.Roles_Undefined
}

func (user *BaseUser) GetStatus() LiveStatus {
	return user.Status
}

func (user *BaseUser) SetStatus(status LiveStatus) {
	user.Status = status
}

func (user *BaseUser) MakeNightMove(proto.MafiaServiceClient) error {
	return fmt.Errorf("not implemented")
}

func MakeUser(login string, role proto.Roles) User {
	base := BaseUser{
		Status: Alive,
		Login:  login,
	}
	if role == proto.Roles_Civilian {
		return &Civilian{
			BaseUser: base,
		}
	} else if role == proto.Roles_Mafia {
		return &Mafia{
			BaseUser: base,
		}
	}
	return &base
}
