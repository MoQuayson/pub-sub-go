package grpc

import (
	"context"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

type PublisherGrpcService struct {
	Id     string
	Conn   *grpc.ClientConn
	Client pb.BrokerClient
	Config *models.PublisherConfig
}

func (p *PublisherGrpcService) Close() error {
	return p.Conn.Close()
}

// PublishMessage publishes data to a topic via a broker
func (p *PublisherGrpcService) PublishMessage(ctx context.Context, topic string, partition models.Partition, data string) error {
	msgId, err := utils.NewMessageId()
	if err != nil {
		return err
	}
	_, err = p.Client.Publish(ctx, &pb.PublishRequest{
		Id:        msgId,
		Topic:     topic,
		Partition: int32(partition),
		Data:      data,
		Timestamp: time.Now().Format(time.DateTime),
	})
	if err != nil {
		slog.Error("publisher failed to publish message", "publisherId", p.Id, "error", err)
		return err
	}
	return nil
}
