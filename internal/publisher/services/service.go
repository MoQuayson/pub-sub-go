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

type PublisherService struct {
	client *rpc.Client
}

func NewPublisherService(cfg *models.RpcConnConfig) *PublisherService {
	pub := &PublisherService{}
	client, err := connectToRpcServer(cfg)

	if err != nil {
		log.Printf("failed to create a publisher: %v\n", err)
		return nil
	}

	pub.client = client
	return pub
}

// PublishMessage publishes to a topic
func (p *PublisherService) PublishMessage(topic string, partition models.Partition, data string) error {
	msg := models.Message{
		Id:        utils.NewMessageId(),
		Topic:     topic,
		Partition: partition,
		Data:      data,
		Timestamp: time.Now(),
	}

	//make an rpc call to a service method
	var reply string
	if err := p.client.Call(constants.PublishToBrokerServiceMethod, msg, &reply); err != nil {
		return err
	}

	return nil

}

func connectToRpcServer(c *models.RpcConnConfig) (*rpc.Client, error) {
	if c == nil {
		*c = models.RpcConnConfig{}
	}
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
}
