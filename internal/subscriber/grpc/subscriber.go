package grpc

import (
	"context"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

type SubscriberGrpcService struct {
	Id     string
	Conn   *grpc.ClientConn
	Client pb.BrokerClient
	Config *models.SubscriberConfig
}

func (s *SubscriberGrpcService) Close() error {
	return s.Conn.Close()
}

func (s *SubscriberGrpcService) Subscribe(ctx context.Context, topic string) ([]*models.Message, error) {
	req := &pb.GetMessagesRequest{
		SubscriberId: s.Id,
		Topic:        topic,
		Partition:    int32(s.Config.Partition),
		PublishTime:  pb.PublishTime(s.Config.MessagePublishTime),
	}

	resp, err := s.Client.GetMessages(ctx, req)
	if err != nil {
		slog.Error("failed to get messages from broker", "error", err)
		return nil, err
	}

	messages := make([]*models.Message, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		messages = append(messages, &models.Message{
			Id:        m.GetMessageId(),
			Topic:     m.GetTopic(),
			Partition: models.Partition(m.GetPartition()),
			Data:      m.GetData(),
			Timestamp: time.Unix(m.GetTimestamp(), 0),
		})
	}
	return messages, nil
}
