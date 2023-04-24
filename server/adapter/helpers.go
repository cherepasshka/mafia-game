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
	defer func() {
		// todo:
		// is this really necessary?
		// delete(s.channels, req.Login)
	}()
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
