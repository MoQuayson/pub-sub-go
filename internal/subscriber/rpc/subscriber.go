package rpc

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"net/rpc"
)

type SubscriberRpcService struct {
	Id     string
	Client *rpc.Client
	Config *models.SubscriberConfig
}

func (s *SubscriberRpcService) getMessages(topic string, partition models.Partition) (models.MessageList, error) {
	request := models.GetMessageRequest{
		SubscriberId: s.Id,
		Topic:        topic,
		Partition:    partition,
		PublishTime:  s.Config.MessagePublishTime,
	}

	var messages models.MessageList
	err := s.Client.Call(constants.GetMessagesFromBrokerServiceMethod, request, &messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *SubscriberRpcService) Subscribe(topic string) ([]*models.Message, error) {
	return s.getMessages(topic, s.Config.Partition)
}

func connectToRpcServer(c *models.SubscriberConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.SubscriberConfig{
			Host: constants.DefaultHost,
			Port: constants.DefaultPort,
		}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
