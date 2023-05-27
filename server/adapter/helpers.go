package mafia_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	segkafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"soa.mafia-game/kafka-help"
	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-game"
	"soa.mafia-game/server/domain/models/party"
)

func (adapter *ServerAdapter) SendReadinessNotification(members []string) {
	adapter.conn_guard.Lock()
	defer adapter.conn_guard.Unlock()
	for _, member := range members {
		channel, exist := adapter.connections[member]
		if exist {
			channel <- mafia_domain.Event{SessionReadiness: true}
		}
	}
}

func (adapter *ServerAdapter) HandleUserMessages(ctx context.Context, user_login string) {
	brokers := strings.Split(adapter.brokerServers, ",")

	reader := segkafka.NewReader(segkafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     user_login,
		Partition: 0,
	})
	reader.SetOffset(0)
	producer, _ := kafka_service.GetNewProducer(adapter.brokerServers)
	groupSession := adapter.getPartySessionId(user_login)
	log.Printf("START LISTENING %v\n", user_login)
	admin, _ := kafka.NewAdminClientFromProducer(producer)
	kafka_service.CreateTopic(admin, groupSession, party.PARTY_SIZE) // TODO
	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(context.Canceled, err) {
				return
			} else {
				log.Printf("%v\n", err)
			}
		}
		fmt.Printf("GOT MESSAGE %v: %v\n", message.Key, message.Value)
		for i := 0; i < party.PARTY_SIZE; i++ {
			err = kafka_service.Produce(user_login, string(message.Value), groupSession, int32(i), producer)
			if err != nil {
				log.Printf("GOT ERROR: %v\n", err)
			}
		}
		// select {
		// 	case <-ctx.Done():
		// 		return
		// }
		// log.Printf("After waiting\n")
	}
}

func (adapter *ServerAdapter) ConnectToSession(ctx context.Context, req *proto.DefaultRequest) (*proto.ConnectToSessionResponse, error) {
	success, event := adapter.game.AddPlayer(req.Login)
	response := &proto.ConnectToSessionResponse{
		Success: success,
		Readiness: &proto.SessionReadiness{
			SessionReady: false,
			Role:         proto.Roles_Undefined,
			SessionId:    adapter.getPartySessionId(req.Login),
		},
	}
	if !success {
		return response, nil
	}
	// create topics
	producer, _ := kafka_service.GetNewProducer(adapter.brokerServers)
	admin, _ := kafka.NewAdminClientFromProducer(producer)
	err := kafka_service.CreateTopic(admin, req.Login, 1)
	if err != nil {
		adapter.game.RemovePlayer(req.Login)
		return &proto.ConnectToSessionResponse{Success: false}, err
	}
	// kafka_service.Produce("aba", req.Login, req.Login, 0, producer)
	// TODO mb through ctx push canceling and deleting topics
	newCtx, cancel := context.WithCancel(context.Background())
	adapter.callbacks_guard.Lock()
	adapter.user_callbacks[req.Login] = cancel
	adapter.callbacks_guard.Unlock()
	go adapter.HandleUserMessages(newCtx, req.Login)

	adapter.conn_guard.Lock()
	defer adapter.conn_guard.Unlock()
	for _, channel := range adapter.connections {
		channel <- event
	}
	fmt.Printf("Hi %v!\n", req.Login)
	adapter.victims[req.Login] = make(chan string, 1)
	return response, nil
}

func (adapter *ServerAdapter) getPartySessionId(user_login string) string {
	return fmt.Sprintf("chat-%v", adapter.game.GetParty(user_login))
}

func (adapter *ServerAdapter) closeConnection(user_login string) {
	adapter.conn_guard.Lock()
	defer adapter.conn_guard.Unlock()
	close(adapter.connections[user_login])
	delete(adapter.connections, user_login)
}

func (adapter *ServerAdapter) CloseChannels(user_login string) {
	adapter.guard.Lock()
	defer adapter.guard.Unlock()
	_, exist := adapter.connections[user_login]
	if exist {
		adapter.closeConnection(user_login)
	}
	start_next_day, exist := adapter.victims[user_login]
	if exist {
		close(start_next_day)
		delete(adapter.victims, user_login)
	}
}

func (adapter *ServerAdapter) LeaveSession(ctx context.Context, request *proto.DefaultRequest) (*proto.LeaveSessionResponse, error) {
	success, event := adapter.game.RemovePlayer(request.Login)
	adapter.conn_guard.Lock()
	for _, channel := range adapter.connections {
		channel <- event
	}
	adapter.conn_guard.Unlock()

	producer, _ := kafka_service.GetNewProducer(adapter.brokerServers)
	admin, _ := kafka.NewAdminClientFromProducer(producer)
	kafka_service.DeleteTopic(admin, request.Login)

	fmt.Printf("Bye %v!\n", request.Login)
	adapter.CloseChannels(request.Login)
	return &proto.LeaveSessionResponse{Success: success}, nil
}

func (adapter *ServerAdapter) ListConnections(req *proto.DefaultRequest, stream proto.MafiaService_ListConnectionsServer) error {
	adapter.game.EnterSession(req.Login)
	_, exist := adapter.connections[req.Login]
	if exist {
		adapter.closeConnection(req.Login)
	}
	msgChannel := make(chan mafia_domain.Event, len(adapter.game.Events)+1)
	adapter.conn_guard.Lock()
	adapter.connections[req.Login] = msgChannel
	adapter.conn_guard.Unlock()
	for i := 0; i < len(adapter.game.Events); i++ {
		msgChannel <- adapter.game.Events[i]
	}
	if adapter.game.SessionReady(req.Login) {
		if adapter.game.DistributeRoles(adapter.game.GetParty(req.Login)) {
			partyId := adapter.getPartySessionId(req.Login)
			producer, _ := kafka_service.GetNewProducer(adapter.brokerServers)
			admin, _ := kafka.NewAdminClientFromProducer(producer)
			err := kafka_service.CreateTopic(admin, partyId, party.PARTY_SIZE)
			if err != nil {
				adapter.closeConnection(req.Login)
				return err
			}
			adapter.SendReadinessNotification(adapter.game.GetMembers(adapter.game.GetParty(req.Login)))
		}
	}
	for {
		select {
		case <-stream.Context().Done():
			adapter.closeConnection(req.Login)
			return nil
		case msg, success := <-msgChannel:
			if !success {
				return nil
			}
			response := &proto.ListConnectionsResponse{
				Login: msg.User,
				State: msg.Status,
				Time:  timestamppb.New(msg.Time),
				Readiness: &proto.SessionReadiness{
					SessionReady: msg.SessionReadiness,
					Role:         proto.Roles_Undefined,
					SessionId:    adapter.getPartySessionId(req.Login),
					Partition:    int32(adapter.game.GetPartition(req.Login)),
				},
			}
			if msg.SessionReadiness {
				response.Readiness.Role = adapter.game.GetRole(req.Login)
				response.Readiness.Players = adapter.game.GetMembers(adapter.game.GetParty(req.Login))
			}
			err := stream.Send(response)
			if err != nil || response.Readiness.SessionReady {
				adapter.closeConnection(req.Login)
				return err
			}
		}
	}
}

func (adapter *ServerAdapter) MakeMove(ctx context.Context, req *proto.MoveRequest) (*proto.MoveResponse, error) {
	role := adapter.game.GetRole(req.Login)
	party := adapter.game.GetParty(req.Login)
	response := &proto.MoveResponse{}
	adapter.guard.Lock()
	defer adapter.guard.Unlock()
	if role == proto.Roles_Civilian {
		adapter.moved_players[party]++
	} else if role == proto.Roles_Commissioner {
		if adapter.game.GetRole(req.Target) == proto.Roles_Mafia {
			response.Accepted = true
		} else {
			response.Accepted = false
		}
		adapter.moved_players[party]++

	} else if role == proto.Roles_Mafia {
		if adapter.game.IsPlayerAlive(req.Target) {
			adapter.game.RecentVictim[party] = req.Target
			response.Accepted = true
			adapter.moved_players[party]++
		} else {
			response.Accepted = false
		}
	}
	alive_cnt := len(adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)))
	if adapter.moved_players[party] == alive_cnt {
		adapter.game.Kill(adapter.game.RecentVictim[party])
		victim := adapter.game.RecentVictim[party]
		adapter.game.RecentVictim[party] = "None"
		for _, member := range adapter.game.GetMembers(adapter.game.GetParty(req.Login)) {
			adapter.victims[member] <- victim
		}
		adapter.moved_players[party] -= alive_cnt
	}
	return response, nil
}

func (adapter *ServerAdapter) StartDay(ctx context.Context, req *proto.DefaultRequest) (*proto.DayResponse, error) {
	victim := <-adapter.victims[req.Login]
	resp := &proto.DayResponse{
		Victim: victim,
		Alive:  adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)),
	}
	return resp, nil
}

func (adapter *ServerAdapter) VoteForMafia(ctx context.Context, req *proto.VoteForMafiaRequest) (*proto.VoteForMafiaResponse, error) {
	adapter.game.VoteFor(req.Login, req.MafiaGuess)
	ghost := adapter.game.WaitForEverybody(req.Login)
	response := &proto.VoteForMafiaResponse{KilledUser: ghost, KilledUserRole: adapter.game.GetRole(ghost)}
	return response, nil
}

func (adapter *ServerAdapter) GetStatus(ctx context.Context, req *proto.DefaultRequest) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Alive: adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)),
		GameStatus: &proto.GameStatus{
			Active: adapter.game.IsActive(adapter.game.GetParty(req.Login)),
			Winner: adapter.game.Winner(adapter.game.GetParty(req.Login)),
		},
	}, nil
}

func (adapter *ServerAdapter) ExitGameSession(ctx context.Context, req *proto.DefaultRequest) (*proto.ExitGameSessionResponse, error) {
	adapter.callbacks_guard.Lock()
	adapter.user_callbacks[req.Login]()
	adapter.callbacks_guard.Unlock()
	partyId := adapter.getPartySessionId(req.Login)
	producer, _ := kafka_service.GetNewProducer(adapter.brokerServers)
	admin, _ := kafka.NewAdminClientFromProducer(producer)
	kafka_service.DeleteTopic(admin, partyId)
	adapter.game.ExitSession(req.Login)
	return &proto.ExitGameSessionResponse{}, nil
}
