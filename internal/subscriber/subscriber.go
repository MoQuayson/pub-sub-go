package subscriber

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/subscriber"
	"log"
	"net/rpc"
	"time"
)

type SubscriberService struct {
	subscriber.Subscriber
	id     string
	client *rpc.Client
}

func NewSubscriberService(cfg *models.RpcConnConfig) subscriber.Subscriber {
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

func connectToRpcServer(c *models.RpcConnConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.RpcConnConfig{}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
