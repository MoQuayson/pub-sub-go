package services

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/constants"
	"log"
	"net/rpc"
)

type SubscriberService struct {
	id     string
	client *rpc.Client
	config *models.SubscriberConfig
}

func NewSubscriberService(cfg *models.SubscriberConfig) *SubscriberService {
	client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Fatalf("failed to create subscriber: %v\n", err)
		return nil
	}

	subId := utils.NewSubscriberId()
	if cfg != nil && cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &SubscriberService{
		client: client,
		id:     subId,
		config: cfg,
	}
}

func (s *SubscriberService) getMessages(topic string, partition models.Partition) (models.MessageList, error) {
	request := models.GetMessageRequest{
		SubscriberId: s.id,
		Topic:        topic,
		Partition:    partition,
		PublishTime:  s.config.MessagePublishTime,
	}

	var messages models.MessageList
	err := s.client.Call(constants.GetMessagesFromBrokerServiceMethod, request, &messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *SubscriberService) Subscribe(topic string) (models.MessageList, error) {
	return s.getMessages(topic, s.config.Partition)
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
