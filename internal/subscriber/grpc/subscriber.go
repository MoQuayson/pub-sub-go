package grpc

import (
	"context"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"time"
)

type SubscriberGrpcService struct {
	Id     string
	Client pb.BrokerClient
	Config *models.SubscriberConfig
}

func (s *SubscriberGrpcService) Subscribe(topic string) ([]*models.Message, error) {
	req := &pb.GetMessagesRequest{
		SubscriberId: s.Id,
		Topic:        topic,
		Partition:    int32(s.Config.Partition),
		PublishTime:  pb.PublishTime(s.Config.MessagePublishTime),
	}

	resp, err := s.Client.GetMessages(context.Background(), req)
	if err != nil {
		log.Printf("failed to get messages from broker: %v\n", err)
		return nil, err
	}

	messages := make([]*models.Message, 0)

	for _, m := range resp.Messages {
		messages = append(messages, &models.Message{
			Topic:     m.GetTopic(),
			Partition: models.Partition(m.GetPartition()),
			Data:      m.GetData(),
			Timestamp: time.Unix(m.GetTimestamp(), 0),
		})
	}

	return messages, nil
}
