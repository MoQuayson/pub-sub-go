package grpc

import (
	"context"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/internal/broker"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"sync"
)

type BrokerGrpcServer struct {
	Config            *models.BrokerConfig
	Storage           storage.Storage
	SubscriberOffsets models.SubscriberOffsets
	Mutex             sync.Mutex
	Server            *Server
	Logger            *slog.Logger
}

func (b *BrokerGrpcServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.Config.Host, b.Config.Port))
	if err != nil {
		return fmt.Errorf("broker: listen on %s:%s: %w", b.Config.Host, b.Config.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBrokerServer(grpcServer, b.Server)
	b.Server.broker = b
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		b.Logger.Info("broker shutting down")
		grpcServer.GracefulStop()
	}()

	b.Logger.Info("broker listening", "host", b.Config.Host, "port", b.Config.Port)
	if err = grpcServer.Serve(lis); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("broker: serve: %w", err)
	}
	return nil
}

// PublishMessage publishes a message to a topic
func (b *BrokerGrpcServer) PublishMessage(msg *models.Message) error {
	if err := b.Storage.StoreMessage(msg); err != nil {
		return err
	}
	return nil
}

func (b *BrokerGrpcServer) GetMessages(req *models.GetMessageRequest) ([]*models.Message, error) {
	messages, err := b.Storage.GetMessages(req.Topic, req.Partition)
	if err != nil {
		b.Logger.Error("failed to get messages", "error", err)
		return nil, err
	}
	return broker.ResolveMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
}
