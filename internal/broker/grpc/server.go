package grpc

import (
	"context"
	"fmt"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"time"
)

type Server struct {
	pb.UnimplementedBrokerServer
	broker *BrokerGrpcServer
}

func (b *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.StringValue, error) {

	timestamp, err := time.Parse(time.DateTime, req.GetTimestamp())
	if err != nil {
		return nil, err
	}
	msg := &models.Message{
		Id:        req.GetId(),
		Topic:     req.GetTopic(),
		Partition: models.Partition(req.GetPartition()),
		Data:      req.GetData(),
		Timestamp: timestamp,
	}
	if err = b.broker.PublishMessage(msg); err != nil {
		return nil, err
	}

	return &pb.StringValue{Value: fmt.Sprintf("message (%s) published successfully", msg.Id)}, nil
}

func (b *Server) StreamMessages(req *pb.GetMessagesRequest, stream pb.Broker_StreamMessagesServer) error {
	messages, err := b.GetMessages(context.Background(), req)
	if err != nil {
		return err
	}

	for _, msg := range messages.Messages {
		if err = stream.Send(msg); err != nil {
			log.Printf("Error sending message to subscriber: %v", err)
			break
		}
	}

	return nil
}
func (b *Server) GetMessages(_ context.Context, req *pb.GetMessagesRequest) (*pb.MessageList, error) {
	messages, err := b.broker.GetMessages(&models.GetMessageRequest{
		SubscriberId: req.GetSubscriberId(),
		Topic:        req.GetTopic(),
		Partition:    models.Partition(req.GetPartition()),
		PublishTime:  models.PublishTime(req.GetPublishTime()),
	})

	if err != nil {
		log.Printf("failed to get messages from broker: %v\n", err)
		return nil, err
	}

	pbMessages := make([]*pb.Message, 0)

	for _, message := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId: message.Id,
			Topic:     message.Topic,
			Partition: int32(message.Partition),
			Data:      message.Data,
			Timestamp: message.Timestamp.Unix(),
		})
	}

	return &pb.MessageList{Messages: pbMessages}, err
}
