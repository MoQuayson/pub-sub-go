package subscriber

import (
	"fmt"
	models "github.com/MoQuayson/go-event-bridge/pkg/shared/models"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils"
	"github.com/MoQuayson/go-event-bridge/pkg/shared/utils/constants"
	"log"
	"net/rpc"
	"time"
)

type Subscriber struct {
	id     string
	client *rpc.Client
}

func NewSubscriber(cfg *models.RpcConnConfig) *Subscriber {
	client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Printf("failed to create a subscriber: %v\n", err)
		return nil
	}

	return &Subscriber{
		client: client,
		id:     utils.NewSubscriberId(),
	}
}

func (s *Subscriber) GetMessages(topic string, partition models.Partition, startTime time.Time) (models.MessageList, error) {
	//request := map[string]interface{}{
	//	"subscriberID": subscriberID,
	//	"topic":        topic,
	//	"partition":    partition,
	//	"startTime":    startTime,
	//}

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
