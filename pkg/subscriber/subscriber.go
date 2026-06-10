package subscriber

import (
	"context"
	"fmt"
	subGrpc "github.com/MoQuayson/pub-sub-go/internal/subscriber/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/subscriber/rpc"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) ([]*models.Message, error)
}

func NewSubscriber(cfg *models.SubscriberConfig) (Subscriber, error) {
	if cfg == nil {
		return nil, fmt.Errorf("subscriber: config is required")
	}

	switch cfg.Transport {
	case "", models.DefaultTransport:
		return initRpcSubscriber(cfg)
	case models.GrpcTransport:
		return initGrpcSubscriber(cfg)
	default:
		return nil, fmt.Errorf("subscriber: unsupported transport %q", cfg.Transport)
	}
}

func initRpcSubscriber(cfg *models.SubscriberConfig) (Subscriber, error) {
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("subscriber: connect to rpc broker: %w", err)
	}

	subId, err := utils.NewSubscriberId()
	if err != nil {
		return nil, fmt.Errorf("subscriber: generate id: %w", err)
	}
	if cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &rpc.SubscriberRpcService{
		Client: client,
		Id:     subId,
		Config: cfg,
	}, nil
}

func initGrpcSubscriber(cfg *models.SubscriberConfig) (Subscriber, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("subscriber: connect to grpc broker: %w", err)
	}

	subId, err := utils.NewSubscriberId()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("subscriber: generate id: %w", err)
	}
	if cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &subGrpc.SubscriberGrpcService{
		Conn:   conn,
		Client: pb.NewBrokerClient(conn),
		Id:     subId,
		Config: cfg,
	}, nil
}
