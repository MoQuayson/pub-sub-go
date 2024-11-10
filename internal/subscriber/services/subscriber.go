package services

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/constants"
	"log"
	"net/rpc"
	"time"
)

type SubscriberService struct {
	id     string
	client *rpc.Client
}

func NewSubscriberService(cfg *models.BrokerConfig) *SubscriberService {
	client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Printf("failed to create a subscriber: %v\n", err)
		return nil
	}

	return &SubscriberService{
		client: client,
		id:     utils.NewSubscriberId(),
	}
}

func (s *SubscriberService) GetMessages(topic string, partition models.Partition, startTime time.Time) (models.MessageList, error) {
	request := models.GetMessageRequest{
		SubscriberId: s.id,
		Topic:        topic,
		Partition:    partition,
		Timestamp:    startTime,
	}

	var messages models.MessageList
	err := s.client.Call(constants.GetMessagesFromBrokerServiceMethod, request, &messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func connectToRpcServer(c *models.BrokerConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.BrokerConfig{}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
