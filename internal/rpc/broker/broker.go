package broker

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	linq "github.com/samber/lo"
	"log"
	"net"
	"net/rpc"
	"sort"
	"sync"
	"time"
)

type RpcBroker struct {
	Config            *models.BrokerConfig
	Storage           storage.Storage
	SubscriberOffsets models.SubscriberOffsets
	Mutex             sync.Mutex
	//server            server.Server
}

func (b *RpcBroker) Start() {
	//register broker first
	if err := b.Register(b); err != nil {
		log.Fatalf("failed to register broker on rpc: %v\n", err)
	}

	//run on tcp
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.Config.Host, b.Config.Port))
	if err != nil {
		log.Println("failed to start broker:", err)
		return
	}
	defer listener.Close()
	log.Printf("Broker is running on %s:%s", b.Config.Host, b.Config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// Publish stores message in a storage
func (b *RpcBroker) Publish(msg *models.Message, reply *string) error {
	if err := b.Storage.StoreMessage(msg); err != nil {
		return err
	}

	*reply = fmt.Sprintf("message (%s) published successfully", msg.Id)
	return nil
}

// GetMessages retrieves messages for subscriber
func (b *RpcBroker) GetMessages(request *models.GetMessageRequest, reply *models.MessageList) error {
	var err error
	switch request.PublishTime {
	case models.LatestPublishTime:
		*reply, err = b.getLatestMessages(request)
		return err
	case models.EarliestPublishTime:
		*reply, err = b.getEarliestMessages(request)
		return err
	default:
		*reply, err = b.getMessages(request, time.Duration(request.PublishTime))
		return err
	}
}

func (b *RpcBroker) getEarliestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
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

func (b *RpcBroker) getLatestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
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

func (b *RpcBroker) getMessages(request *models.GetMessageRequest, duration time.Duration) (models.MessageList, error) {
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

func (s *RpcBroker) Register(service any) error {
	if err := rpc.Register(service); err != nil {
		return err
	}

	return nil
}
