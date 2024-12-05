package subscriber

import (
	"fmt"
	subGrpc "github.com/MoQuayson/pub-sub-go/internal/subscriber/grpc"
	"github.com/MoQuayson/pub-sub-go/internal/subscriber/rpc"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"log"
)

type Subscriber interface {
	Subscribe(topic string) ([]*models.Message, error)
}

func NewSubscriber(cfg *models.SubscriberConfig) Subscriber {

	if cfg == nil {
		return nil
	}

	switch cfg.Transport {
	case models.GrpcTransport:
		return initGrpcSubscriber(cfg)
	default:
		return initRpcSubscriber(cfg)
	}
}

func initRpcSubscriber(cfg *models.SubscriberConfig) Subscriber {
	client, err := models.ConnectToRpcServer(cfg.Host, cfg.Port)

	if err != nil {
		log.Fatalf("failed to create subscriber: %v\n", err)
		return nil
	}

	subId := utils.NewSubscriberId()
	if cfg != nil && cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	return &rpc.SubscriberRpcService{
		Client: client,
		Id:     subId,
		Config: cfg,
	}
}

func initGrpcSubscriber(cfg *models.SubscriberConfig) Subscriber {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to grpc broker: %v\n", err)
	}
	//defer conn.Close()

	subId := utils.NewSubscriberId()
	if cfg != nil && cfg.SubscriberId != nil {
		subId = *cfg.SubscriberId
	}

	sub := &subGrpc.SubscriberGrpcService{}
	sub.Client = pb.NewBrokerClient(conn)
	sub.Id = subId
	sub.Config = cfg

	return sub
}
