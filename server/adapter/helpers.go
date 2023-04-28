package mafia_server

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-game"
)

func (adapter *ServerAdapter) SendReadinessNotification(members []string) {
	for _, member := range members {
		channel, exist := adapter.connections[member]
		if exist {
			channel <- mafia_domain.Event{SessionReadiness: true}
		}
	}
}

func (adapter *ServerAdapter) ConnectToSession(ctx context.Context, user *proto.User) (*proto.ConnectToSessionResponse, error) {
	success, event := adapter.game.AddPlayer(user.Name)
	response := &proto.ConnectToSessionResponse{
		Success: success,
		Readiness: &proto.SessionReadiness{
			SessionReady: false,
			Role:         proto.Roles_Undefined,
		},
	}
	if !success {
		return response, nil
	}
	for _, channel := range adapter.connections {
		channel <- event
	}
	fmt.Printf("Hi %v!\n", user.Name)
	if adapter.game.SessionReady(user.Name) {
		if adapter.game.DistributeRoles(adapter.game.GetParty(user.Name)) {
			response.Readiness.SessionReady = true
			response.Readiness.Role = adapter.game.GetRole(user.Name)
			response.Readiness.Players = adapter.game.GetMembers(adapter.game.GetParty(user.Name))
			adapter.SendReadinessNotification(adapter.game.GetMembers(adapter.game.GetParty(user.Name)))
		}
	}
	adapter.victims[user.Name] = make(chan string, 1)
	return response, nil
}

func (adapter *ServerAdapter) CloseChannels(user_login string) {
	conection, exist := adapter.connections[user_login]
	if exist {
		close(conection)
		delete(adapter.connections, user_login)
	}
	start_next_day, exist := adapter.victims[user_login]
	if exist {
		close(start_next_day)
		delete(adapter.victims, user_login)
	}
}

func (adapter *ServerAdapter) LeaveSession(ctx context.Context, request *proto.LeaveSessionRequest) (*proto.LeaveSessionResponse, error) {
	success, event := adapter.game.RemovePlayer(request.User.Name)
	for _, channel := range adapter.connections {
		channel <- event
	}
	fmt.Printf("Bye %v!\n", request.User.Name)
	adapter.CloseChannels(request.User.Name)
	return &proto.LeaveSessionResponse{Success: success}, nil
}

func (adapter *ServerAdapter) ListConnections(req *proto.ListConnectionsRequest, stream proto.MafiaService_ListConnectionsServer) error {
	msgChannel := make(chan mafia_domain.Event, len(adapter.game.Events)+1)
	adapter.connections[req.Login] = msgChannel
	for i := 0; i < len(adapter.game.Events); i++ {
		msgChannel <- adapter.game.Events[i]
	}
	for {
		select {
		case <-stream.Context().Done():
			close(msgChannel)
			delete(adapter.connections, req.Login)
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
				},
			}
			if msg.SessionReadiness {
				response.Readiness.Role = adapter.game.GetRole(req.Login)
				response.Readiness.Players = adapter.game.GetMembers(adapter.game.GetParty(req.Login))
			}
			err := stream.Send(response)
			if err != nil {
				close(msgChannel)
				return err
			}
			if response.Readiness.SessionReady {
				return nil
			}
		}
	}
}

func (adapter *ServerAdapter) MakeMove(ctx context.Context, req *proto.MoveRequest) (*proto.MoveResponse, error) {
	role := adapter.game.GetRole(req.Login)
	party := adapter.game.GetParty(req.Login)
	response := &proto.MoveResponse{}
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
	adapter.guard.Lock()
	defer adapter.guard.Unlock()
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

func (adapter *ServerAdapter) StartDay(ctx context.Context, req *proto.DayRequest) (*proto.DayResponse, error) {
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

func (adapter *ServerAdapter) GetStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Alive: adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)),
		GameStatus: &proto.GameStatus{
			Active: adapter.game.IsActive(adapter.game.GetParty(req.Login)),
			Winner: adapter.game.Winner(adapter.game.GetParty(req.Login)),
		},
	}, nil
}
