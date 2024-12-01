package subscriber

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	constants2 "github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"net/rpc"
)

type RpcSubscriber struct {
	Id     string
	Client *rpc.Client
	Config *models.SubscriberConfig
}

func NewSubscriberService(cfg *models.SubscriberConfig) *RpcSubscriber {
	Client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Fatalf("failed to create subscriber: %v\n", err)
		return nil
	}

	subId := utils.NewSubscriberId()
	if cfg != nil && cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &RpcSubscriber{
		Client: Client,
		Id:     subId,
		Config: cfg,
	}
}

func (s *RpcSubscriber) getMessages(topic string, partition models.Partition) (models.MessageList, error) {
	request := models.GetMessageRequest{
		SubscriberId: s.Id,
		Topic:        topic,
		Partition:    partition,
		PublishTime:  s.Config.MessagePublishTime,
	}

	var messages models.MessageList
	err := s.Client.Call(constants2.GetMessagesFromBrokerServiceMethod, request, &messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *RpcSubscriber) Subscribe(topic string) (models.MessageList, error) {
	return s.getMessages(topic, s.Config.Partition)
}

func connectToRpcServer(c *models.SubscriberConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.SubscriberConfig{
			Host: constants2.DefaultHost,
			Port: constants2.DefaultPort,
		}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
