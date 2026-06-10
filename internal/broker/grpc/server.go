package grpc

import (
	"context"
	"fmt"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Server struct {
	pb.UnimplementedBrokerServer
	broker *BrokerGrpcServer
}

func (b *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.StringValue, error) {
	timestamp, err := time.Parse(time.DateTime, req.GetTimestamp())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid timestamp %q: %v", req.GetTimestamp(), err)
	}

	msg := &models.Message{
		Id:        req.GetId(),
		Topic:     req.GetTopic(),
		Partition: models.Partition(req.GetPartition()),
		Data:      req.GetData(),
		Timestamp: timestamp,
	}
	if err = b.broker.PublishMessage(msg); err != nil {
		return nil, status.Errorf(codes.Internal, "publish failed: %v", err)
	}

	return &pb.StringValue{Value: fmt.Sprintf("message (%s) published successfully", msg.Id)}, nil
}

func (b *Server) StreamMessages(req *pb.GetMessagesRequest, stream pb.Broker_StreamMessagesServer) error {
	messages, err := b.GetMessages(stream.Context(), req)
	if err != nil {
		return err
	}

	for _, msg := range messages.Messages {
		if err = stream.Send(msg); err != nil {
			return status.Errorf(codes.Internal, "stream send failed: %v", err)
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
		return nil, status.Errorf(codes.Internal, "get messages failed: %v", err)
	}

	pbMessages := make([]*pb.Message, 0, len(messages))
	for _, message := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId: message.Id,
			Topic:     message.Topic,
			Partition: int32(message.Partition),
			Data:      message.Data,
			Timestamp: message.Timestamp.Unix(),
		})
	}

	return &pb.MessageList{Messages: pbMessages}, nil
}
