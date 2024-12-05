package grpc

import (
	"context"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"time"
)

type PublisherGrpcService struct {
	Id     string
	Client pb.BrokerClient
	Config *models.PublisherConfig
}

// PublishMessage publishes a data to a topic via a broker
func (p *PublisherGrpcService) PublishMessage(topic string, partition models.Partition, data string) error {
	_, err := p.Client.Publish(context.Background(), &pb.PublishRequest{
		Id:        utils.NewMessageId(),
		Topic:     topic,
		Partition: int32(partition),
		Data:      data,
		Timestamp: time.Now().Format(time.DateTime),
	})
	if err != nil {
		log.Printf("publisher (%s) failed to publish message: %v\n", p.Id, err)
		return err
	}
	return nil
}
