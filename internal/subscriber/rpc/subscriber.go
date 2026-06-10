package rpc

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"net/rpc"
)

type SubscriberRpcService struct {
	Id     string
	Client *rpc.Client
	Config *models.SubscriberConfig
}

// Subscribe retrieves messages from the broker.
// TODO: propagate ctx to RPC call when net/rpc supports it
func (s *SubscriberRpcService) Subscribe(_ context.Context, topic string) ([]*models.Message, error) {
	request := models.GetMessageRequest{
		SubscriberId: s.Id,
		Topic:        topic,
		Partition:    s.Config.Partition,
		PublishTime:  s.Config.MessagePublishTime,
	}

	var messages models.MessageList
	err := s.Client.Call(constants.GetMessagesFromBrokerServiceMethod, request, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
