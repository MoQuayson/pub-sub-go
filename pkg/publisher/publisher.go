package publisher

import (
	"fmt"
	pubGrpc "github.com/MoQuayson/pub-sub-go/internal/publisher/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/publisher/rpc"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"log"
)

type Publisher interface {
	PublishMessage(topic string, partition models.Partition, data string) error
}

func NewPublisher(cfg *models.PublisherConfig) Publisher {

	if cfg == nil {
		return nil
	}

	switch cfg.Transport {
	case models.GrpcTransport:
		return initGrpcPublisher(cfg)
	default:
		return initRpcPublisher(cfg)
	}
}

func initRpcPublisher(cfg *models.PublisherConfig) Publisher {
	pub := &rpc.PublisherRpcService{}
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)

	if err != nil {
		log.Fatalf("failed to create a publisher: %v\n", err)
		return nil
	}

	pub.Client = client
	return pub
}

func initGrpcPublisher(cfg *models.PublisherConfig) Publisher {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to grpc broker: %v\n", err)
	}
	//defer conn.Close()

	pub := &pubGrpc.PublisherGrpcService{}
	pub.Client = pb.NewBrokerClient(conn)
	pub.Id = utils.NewPublisherId()
	pub.Config = cfg

	return pub
}
