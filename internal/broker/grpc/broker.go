package grpc

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/internal/broker"
	pb "github.com/MoQuayson/pub-sub-go/pkg/proto_gen/github.com/moquayson/pub-sub-go"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	linq "github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"sync"
	"time"
)

type BrokerGrpcServer struct {
	Config            *models.BrokerConfig
	Storage           storage.Storage
	SubscriberOffsets models.SubscriberOffsets
	Mutex             sync.Mutex
	Server            *Server
}

func (b *BrokerGrpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.Config.Host, b.Config.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBrokerServer(grpcServer, b.Server)
	b.Server.broker = b
	// enable reflection
	reflection.Register(grpcServer)
	log.Printf("broker is running on %s:%s\n", b.Config.Host, b.Config.Port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
		return err
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
		log.Printf("failed to get messages: %v\n", err)
		return nil, err
	}
	switch req.PublishTime {
	case models.LatestPublishTime:
		return broker.GetLatestMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
	case models.EarliestPublishTime:
		//return b.getEarliestMessages(req)
		return broker.GetEarliestMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
	default:
		//return b.getMessages(req, time.Duration(req.PublishTime))
		return broker.GetMessages(&b.Mutex, b.SubscriberOffsets, messages, req, time.Duration(req.PublishTime))
	}
}

func (b *BrokerGrpcServer) getEarliestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
	b.Mutex.Lock()
	offset, exists := b.SubscriberOffsets[request.SubscriberId]
	b.Mutex.Unlock()
	//get messages from storage
	messages, err := b.Storage.GetMessages(request.Topic, request.Partition)
	if err != nil {
		return nil, err
	}

	if exists {
		msgCount := len(messages)
		//sort by earliest
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Before(offset.Timestamp) || msg.Timestamp.Equal(offset.Timestamp)
		})

		if len(messages) == msgCount {
			messages = nil
		}
	} else {
		//sort by earliest
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Timestamp.Before(messages[j].Timestamp) || messages[i].Timestamp.Equal(messages[j].Timestamp)
		})
	}

	if len(messages) > 0 {
		b.Mutex.Lock()

		b.SubscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.Mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func (b *BrokerGrpcServer) getLatestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
	b.Mutex.Lock()
	offset, exists := b.SubscriberOffsets[request.SubscriberId]
	b.Mutex.Unlock()
	//get messages from storage
	messages, err := b.Storage.GetMessages(request.Topic, request.Partition)
	if err != nil {
		return nil, err
	}

	if exists {
		//sort by latest
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.After(offset.Timestamp)
		})
	} else {
		//sort by latest
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Timestamp.After(messages[j].Timestamp)
		})
	}
	if len(messages) > 0 {
		b.Mutex.Lock()

		b.SubscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.Mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func (b *BrokerGrpcServer) getMessages(request *models.GetMessageRequest, duration time.Duration) (models.MessageList, error) {
	b.Mutex.Lock()
	offset, exists := b.SubscriberOffsets[request.SubscriberId]
	b.Mutex.Unlock()
	//get messages from storage
	messages, err := b.Storage.GetMessages(request.Topic, request.Partition)
	if err != nil {
		return nil, err
	}

	if duration < 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Sub(time.Now().UTC().Add(duration)) >= 0
		})
	} else {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {

			return time.Now().Sub(msg.Timestamp.Add(duration)) <= 0
		})
	}

	if exists && duration < 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.After(offset.Timestamp)
		})
	} else if exists && duration > 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Before(offset.Timestamp)
		})
	}
	if len(messages) > 0 {
		b.Mutex.Lock()

		b.SubscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.Mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}
