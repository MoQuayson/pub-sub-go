package broker

import (
	"context"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/internal/broker/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/broker/rpc"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log/slog"
)

// Broker interface for all brokers created (rpc,grpc etc)
type Broker interface {
	Start(ctx context.Context) error
	// PublishMessage publishes a message to a topic
	PublishMessage(msg *models.Message) error
	GetMessages(request *models.GetMessageRequest) ([]*models.Message, error)
}

func NewBroker(config *models.BrokerConfig) (Broker, error) {
	if config == nil {
		return nil, fmt.Errorf("broker: config is required")
	}

	switch config.Transport {
	case "", models.DefaultTransport:
		return initRpcBroker(config)
	case models.GrpcTransport:
		return initGrpcBroker(config)
	default:
		return nil, fmt.Errorf("broker: unsupported transport %q", config.Transport)
	}
}

// initRpcBroker initializes broker for rpc communication
func initRpcBroker(config *models.BrokerConfig) (Broker, error) {
	s, err := storage.GetStorage(config)
	if err != nil {
		return nil, err
	}
	b := &rpc.BrokerRpcServer{}
	b.Config = config
	b.Storage = s
	b.SubscriberOffsets = make(map[string]models.Offset)
	b.Logger = slog.Default()
	b.Server = rpc.NewRpcServer(b)
	return b, nil
}

// initGrpcBroker initializes broker for grpc communication
func initGrpcBroker(cfg *models.BrokerConfig) (Broker, error) {
	s, err := storage.GetStorage(cfg)
	if err != nil {
		return nil, err
	}
	return &grpc.BrokerGrpcServer{
		Config:            cfg,
		Storage:           s,
		SubscriberOffsets: make(models.SubscriberOffsets),
		Server:            &grpc.Server{},
		Logger:            slog.Default(),
	}, nil
}
