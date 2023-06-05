package models

import (
	"context"
	"fmt"

	kafka_service "soa.mafia-game/kafka-help"
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
	MakeNightMove(context.Context, []string, proto.MafiaServiceClient) error
	VoteForMafia(context.Context, []string, proto.MafiaServiceClient) error
	Stop()
}

type BaseUser struct {
	Login     string
	Status    LiveStatus
	Session   string
	Partition int32
}

func (user *BaseUser) ExcludeFromAliveList(alive_users []string) []string {
	result := make([]string, 0)
	for i := range alive_users {
		if alive_users[i] != user.Login {
			result = append(result, alive_users[i])
		}
	}
	return result
}

type CommunicatorUser struct {
	BaseUser
	ExitedChat bool
}

func (user *CommunicatorUser) GetLogin() string {
	return user.Login
}

func (user *CommunicatorUser) GetRole() proto.Roles {
	return proto.Roles_Undefined
}

func (user *CommunicatorUser) GetStatus() LiveStatus {
	return user.Status
}

func (user *CommunicatorUser) SetStatus(status LiveStatus) {
	user.Status = status
}

func (user *CommunicatorUser) Stop() {
	if !user.ExitedChat {
		producer, _ := kafka_service.GetNewProducer("localhost:9092")
		defer producer.Close()
		kafka_service.Produce(user.Session, "exit", user.Login, 0, producer)
	}
}

func (user *CommunicatorUser) MakeNightMove(context.Context, []string, proto.MafiaServiceClient) error {
	return fmt.Errorf("not implemented")
}

func (user *CommunicatorUser) VoteForMafia(context.Context, []string, proto.MafiaServiceClient) error {
	return fmt.Errorf("not implemented")
}
