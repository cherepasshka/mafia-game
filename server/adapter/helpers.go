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

func (s *mafiaServer) ConnectToSession(ctx context.Context, user *proto.User) (*proto.ConnectToSessionResponse, error) {
	success, event := s.game.AddPlayer(user.Name)
	for _, channel := range s.channels {
		channel <- event
	}
	fmt.Printf("Hi %v!\n", user.Name)
	return &proto.ConnectToSessionResponse{Success: success}, nil
}

func (s *mafiaServer) LeaveSession(ctx context.Context, request *proto.LeaveSessionRequest) (*proto.LeaveSessionResponse, error) {
	success, event := s.game.RemovePlayer(request.User.Name)
	for _, channel := range s.channels {
		channel <- event
	}
	fmt.Printf("Bye %v!\n", request.User.Name)
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
		delete(s.channels, req.Login)
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
			err := stream.Send(&proto.Connection{Login: msg.User, State: msg.Status, Time: timestamppb.New(msg.Time)})
			if err != nil {
				close(msgChannel)
				return err
			}
		}
	}
}
