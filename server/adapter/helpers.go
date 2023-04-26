package mafia_server

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-game"
)

func (adapter *ServerAdapter) SendNotification(members []string) {
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
	for nm, channel := range adapter.connections {
		fmt.Printf("sent to %s\n", nm)
		channel <- event
	}
	fmt.Printf("Hi %v!\n", user.Name)
	if adapter.game.SessionReady(user.Name) {
		if adapter.game.DistributeRoles(adapter.game.GetParty(user.Name)) {
			response.Readiness.SessionReady = true
			response.Readiness.Role = adapter.game.GetRole(user.Name)
			response.Readiness.Players = adapter.game.GetMembers(adapter.game.GetParty(user.Name))
			adapter.SendNotification(adapter.game.GetMembers(adapter.game.GetParty(user.Name)))
		}
	}
	adapter.start_next_day[user.Name] = make(chan bool, 1)
	return response, nil
}

func (adapter *ServerAdapter) LeaveSession(ctx context.Context, request *proto.LeaveSessionRequest) (*proto.LeaveSessionResponse, error) {
	success, event := adapter.game.RemovePlayer(request.User.Name)
	for _, channel := range adapter.connections {
		channel <- event
	}
	fmt.Printf("Bye %v!\n", request.User.Name)
	// verify that channel is open
	close(adapter.connections[request.User.Name])
	return &proto.LeaveSessionResponse{Success: success}, nil
}

func (adapter *ServerAdapter) ListConnections(req *proto.ListConnectionsRequest, stream proto.MafiaService_ListConnectionsServer) error {
	msgChannel := make(chan mafia_domain.Event, len(adapter.game.Events)+1)
	adapter.connections[req.Login] = msgChannel
	fmt.Printf("open list for %s\n", req.Login)
	for i := 0; i < len(adapter.game.Events); i++ {
		msgChannel <- adapter.game.Events[i]
	}
	for {
		select {
		case <-stream.Context().Done():
			close(msgChannel)
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
			if err != nil || msg.SessionReadiness {
				close(msgChannel)
				return err
			}
		}
	}
}

func (adapter *ServerAdapter) MakeMove(ctx context.Context, req *proto.MoveRequest) (*proto.MoveResponse, error) {
	/*
		mb better to call this function from every player whenewer dead or alive, else мертвяки могут не успеть обновить день и остаться в предыдущем
	*/
	role := adapter.game.GetRole(req.Login)
	response := &proto.MoveResponse{}
	if role == proto.Roles_Civilian {
		adapter.cnt++
	} else if role == proto.Roles_Commissioner {
		if adapter.game.GetRole(req.Target) == proto.Roles_Mafia {
			response.Accepted = true
		} else {
			response.Accepted = false
		}
		adapter.cnt++

	} else if role == proto.Roles_Mafia {
		if adapter.game.IsAlive(req.Target) {
			adapter.game.RecentVictim = req.Target
			response.Accepted = true
			adapter.cnt++
		} else {
			response.Accepted = false
		}
	}
	fmt.Printf("Hi user %s, you are %v, cnt: %v\n", req.Login, role, adapter.cnt)
	adapter.mut.Lock()
	alive_cnt := len(adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)))
	if adapter.cnt == alive_cnt {
		fmt.Printf("In %s send notifications to proceed\n", req.Login)
		adapter.game.Kill(adapter.game.RecentVictim)
		for _, member := range adapter.game.GetMembers(adapter.game.GetParty(req.Login)) {
			fmt.Printf("Sent to %s\n", member)
			adapter.start_next_day[member] <- true
		}
		adapter.cnt -= alive_cnt
	}
	adapter.mut.Unlock()
	return response, nil
}

func (adapter *ServerAdapter) StartDay(ctx context.Context, req *proto.DayRequest) (*proto.DayResponse, error) {
	fmt.Printf("In %s wait to continue and start day\n", req.Login)
	<-adapter.start_next_day[req.Login]
	fmt.Printf("In %s start day\n", req.Login)
	resp := &proto.DayResponse{
		Victim: adapter.game.RecentVictim,
		Alive:  adapter.game.GetAliveMembers(adapter.game.GetParty(req.Login)),
		// GameStatus: &proto.GameStatus{
		// 	Active: adapter.game.IsActive(adapter.game.GetParty(req.Login)),
		// 	Winner: adapter.game.Winner(adapter.game.GetParty(req.Login)),
		// },
	}
	fmt.Printf("for %v alive %v\n", req.Login, resp.Alive)
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
