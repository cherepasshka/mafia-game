package models

import (
	"context"
	"fmt"

	// kafka_service "soa.mafia-game/kafka-help"
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
	// Stop()
}

type BaseUser struct {
	Login     string
	Status    LiveStatus
	Session   string
	Partition int32
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

// func (user *BaseUser) Stop() {
// 	producer, _ := kafka_service.GetNewProducer("localhost:9092")
// 	defer producer.Close()
// 	kafka_service.Produce(user.Session, "exit", user.Login, 0, producer)
// }

func (user *BaseUser) MakeNightMove(context.Context, []string, proto.MafiaServiceClient) error {
	return fmt.Errorf("not implemented")
}

func (user *BaseUser) VoteForMafia(context.Context, []string, proto.MafiaServiceClient) error {
	return fmt.Errorf("not implemented")
}
