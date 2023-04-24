package mafia_server

import (
	"context"
	"fmt"

	// "fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	// "google.golang.org/protobuf/types/known/timestamppb"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_domain "soa.mafia-game/server/domain/mafia-server"
	// timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

func (s *mafiaServer) SendNotification(members []string) {
	for _, member := range members {
		channel, exist := s.channels[member]
		if exist {
			channel <- mafia_domain.Event{SessionReadiness: true}
		}
	}
}

func (s *mafiaServer) ConnectToSession(ctx context.Context, user *proto.User) (*proto.ConnectToSessionResponse, error) {
	success, event := s.game.AddPlayer(user.Name)
	for _, channel := range s.channels {
		channel <- event
	}
	fmt.Printf("Hi %v!\n", user.Name)
	response := &proto.ConnectToSessionResponse{
		Success: success,
		Readiness: &proto.SessionReadiness{
			SessionReady: false,
			Role:         proto.Roles_Undefined,
		},
	}
	if s.game.SessionReady(user.Name) {
		if s.game.DistributeRoles(s.game.GetParty(user.Name)) {
			response.Readiness.SessionReady = true
			response.Readiness.Role = s.game.GetRole(user.Name)
			response.Readiness.Players = s.game.GetMembers(s.game.GetParty(user.Name))
			s.SendNotification(s.game.GetMembers(s.game.GetParty(user.Name)))
		}
	}
	return response, nil
}

func (s *mafiaServer) LeaveSession(ctx context.Context, request *proto.LeaveSessionRequest) (*proto.LeaveSessionResponse, error) {
	success, event := s.game.RemovePlayer(request.User.Name)
	for _, channel := range s.channels {
		channel <- event
	}
	fmt.Printf("Bye %v!\n", request.User.Name)
	// verify that channel is open
	close(s.channels[request.User.Name])
	return &proto.LeaveSessionResponse{Success: success}, nil
}

func (s *mafiaServer) ListConnections(req *proto.ListConnectionsRequest, stream proto.MafiaService_ListConnectionsServer) error {
	msgChannel := make(chan mafia_domain.Event, len(s.game.Events)+1)
	s.channels[req.Login] = msgChannel
	for i := 0; i < len(s.game.Events); i++ {
		msgChannel <- s.game.Events[i]
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
				response.Readiness.Role = s.game.GetRole(req.Login)
				response.Readiness.Players = s.game.GetMembers(s.game.GetParty(req.Login))
			}
			err := stream.Send(response)
			if err != nil || msg.SessionReadiness {
				close(msgChannel)
				return err
			}
		}
	}
}

func (s *mafiaServer) VoteForMafia(context.Context, *proto.VoteForMafiaRequest) (*proto.VoteForMafiaResponse, error) {
	return nil, nil
}

func (s *mafiaServer) MakeMove(ctx context.Context, req *proto.MoveRequest) (*proto.MoveResponse, error) {
	/*
		mb better to call this function from every player whenewer dead or alive, else мертвяки могут не успеть обновить день и остаться в предыдущем
	*/
	role := s.game.GetRole(req.Login)
	response := &proto.MoveResponse{}
	if role == proto.Roles_Civilian {
		s.cnt++
	} else if role == proto.Roles_Commissioner {
		if s.game.GetRole(req.Target) == proto.Roles_Mafia {
			response.Accepted = true
		} else {
			response.Accepted = false
		}
		s.cnt++

	} else if role == proto.Roles_Mafia {
		if s.game.IsAlive(req.Target) {
			s.game.RecentVictim = req.Target
			response.Accepted = true
			s.cnt++
		} else {
			response.Accepted = false
		}
	}
	fmt.Printf("Hi user %s, you are %v, cnt: %v\n", req.Login, role, s.cnt)
	alive_cnt := len(s.game.GetAliveMembers(s.game.GetParty(req.Login)))
	if s.cnt == alive_cnt {
		fmt.Printf("In %s send notifications to proceed\n", req.Login)
		s.game.Kill(s.game.RecentVictim)
		for _, member := range s.game.GetMembers(s.game.GetParty(req.Login)) {
			_, ok := s.ready[member]
			if !ok {
				s.ready[member] = make(chan bool, 1)
			}
			fmt.Printf("Sent to %s\n", member)
			s.ready[member] <- true
		}
		s.cnt -= alive_cnt
	}
	return response, nil
}
func (s *mafiaServer) StartDay(ctx context.Context, req *proto.DayRequest) (*proto.DayResponse, error) {
	fmt.Printf("In %s wait to continue and start day\n", req.Login)
	<-s.ready[req.Login]
	fmt.Printf("In %s start day\n", req.Login)
	resp := &proto.DayResponse{
		Victim: s.game.RecentVictim,
		Alive:  s.game.GetAliveMembers(s.game.GetParty(req.Login)),
	}
	fmt.Printf("for %v alive %v\n", req.Login, resp.Alive)
	return resp, nil
}
