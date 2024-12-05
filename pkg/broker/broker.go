package broker

import (
	"github.com/MoQuayson/pub-sub-go/internal/broker/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/broker/rpc"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
)

// Broker interface for all brokers created (rpc,grpc etc)
type Broker interface {
	Start() error
	//Publish(msg *models.Message, reply *string) error
	// PublishMessage pubishes a message to a topic
	PublishMessage(msg *models.Message) error
	GetMessages(request *models.GetMessageRequest) ([]*models.Message, error)
}

func NewBroker(config *models.BrokerConfig) Broker {

	if config == nil {
		return nil
	}

	switch config.Transport {
	case models.GrpcTransport:
		return initGrpcBroker(config)

	default:
		return initRpcBroker(config)
	}
}

// initRpcBroker initializes broker for rpc communication
func initRpcBroker(config *models.BrokerConfig) Broker {
	b := &rpc.BrokerRpcServer{}
	b.Config = config
	b.Storage = storage.GetStorage(config)
	b.SubscriberOffsets = make(map[string]models.Offset)
	b.Server = rpc.NewRpcServer(b)
	return b
}

// initRpcBroker initializes broker for grpc communication
func initGrpcBroker(cfg *models.BrokerConfig) Broker {
	return &grpc.BrokerGrpcServer{
		Config:            cfg,
		Storage:           storage.GetStorage(cfg),
		SubscriberOffsets: make(models.SubscriberOffsets),
		Server:            &grpc.Server{},
	}
}
