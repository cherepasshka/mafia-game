package chat

type ChatService struct {
	brokerServers string
}

func New(brokerServers string) (*ChatService, error) {
	service := &ChatService{
		brokerServers: brokerServers,
	}
	return service, nil
}
