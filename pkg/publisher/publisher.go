package publisher

import (
	"context"
	"fmt"
	pubGrpc "github.com/MoQuayson/pub-sub-go/internal/publisher/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/publisher/rpc"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Publisher interface {
	PublishMessage(ctx context.Context, topic string, partition models.Partition, data string) error
}

func NewPublisher(cfg *models.PublisherConfig) (Publisher, error) {
	if cfg == nil {
		return nil, fmt.Errorf("publisher: config is required")
	}

	switch cfg.Transport {
	case "", models.DefaultTransport:
		return initRpcPublisher(cfg)
	case models.GrpcTransport:
		return initGrpcPublisher(cfg)
	default:
		return nil, fmt.Errorf("publisher: unsupported transport %q", cfg.Transport)
	}
}

func initRpcPublisher(cfg *models.PublisherConfig) (Publisher, error) {
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("publisher: connect to rpc broker: %w", err)
	}
	return &rpc.PublisherRpcService{Client: client}, nil
}

func initGrpcPublisher(cfg *models.PublisherConfig) (Publisher, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("publisher: connect to grpc broker: %w", err)
	}

	pubId, err := utils.NewPublisherId()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("publisher: generate id: %w", err)
	}

	return &pubGrpc.PublisherGrpcService{
		Conn:   conn,
		Client: pb.NewBrokerClient(conn),
		Id:     pubId,
		Config: cfg,
	}, nil
}
