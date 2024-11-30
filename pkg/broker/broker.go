package broker

import (
	"github.com/MoQuayson/pub-sub-go/internal/rpc/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
)

// Broker interface for all brokers created (rpc,grpc etc)
type Broker interface {
	Start()
	Publish(msg *models.Message, reply *string) error
	GetMessages(request *models.GetMessageRequest, reply *models.MessageList) error
}

func NewBroker(config *models.BrokerConfig) Broker {

	if config == nil {
		return nil
	}

	switch config.Transport {
	default:
		return newRpcBroker(config)
	}
}

func newRpcBroker(config *models.BrokerConfig) Broker {
	b := &broker.RpcBroker{}
	b.Config = config
	b.Storage = storage.GetStorage(config)
	b.SubscriberOffsets = make(map[string]models.Offset)
	return b
}
