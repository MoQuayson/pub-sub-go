package rpc

import (
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"net/rpc"
	"time"
)

type PublisherRpcService struct {
	Client *rpc.Client
}

// PublishMessage publishes to a topic
func (p *PublisherRpcService) PublishMessage(topic string, partition models.Partition, data string) error {
	msg := models.Message{
		Id:        utils.NewMessageId(),
		Topic:     topic,
		Partition: partition,
		Data:      data,
		Timestamp: time.Now(),
	}

	//make an rpc call to a service method
	var reply string
	if err := p.Client.Call(constants.PublishToBrokerServiceMethod, msg, &reply); err != nil {
		return err
	}

	return nil

}
