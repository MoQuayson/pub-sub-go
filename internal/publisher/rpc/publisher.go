package rpc

import (
	"context"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/constants"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"net/rpc"
	"time"
)

type PublisherRpcService struct {
	Client *rpc.Client
}

// PublishMessage publishes to a topic.
// TODO: propagate ctx to RPC call when net/rpc supports it
func (p *PublisherRpcService) PublishMessage(_ context.Context, topic string, partition models.Partition, data string) error {
	msgId, err := utils.NewMessageId()
	if err != nil {
		return err
	}
	msg := models.Message{
		Id:        msgId,
		Topic:     topic,
		Partition: partition,
		Data:      data,
		Timestamp: time.Now(),
	}

	var reply string
	if err := p.Client.Call(constants.PublishToBrokerServiceMethod, msg, &reply); err != nil {
		return err
	}
	return nil
}
