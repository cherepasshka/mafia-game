package mafia_client

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	proto "soa.mafia-game/proto/mafia-game"
)

type Client struct {
	proto.MafiaServiceClient
	shutdownFuncs []func() error
}

func New(host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := &Client{
		MafiaServiceClient: proto.NewMafiaServiceClient(conn),
	}
	client.shutdownFuncs = append(client.shutdownFuncs, conn.Close)
	return client, nil
}

func (c *Client) Stop() {
	for i := len(c.shutdownFuncs) - 1; i >= 0; i-- {
		err := c.shutdownFuncs[i]()
		if err != nil {
			log.Printf("error %v", err)
		}
	}
}
