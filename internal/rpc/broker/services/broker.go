package services

import (
	"encoding/gob"
	"fmt"
	server "github.com/MoQuayson/pub-sub-go/internal/rpc"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/utils/enums"
	linq "github.com/samber/lo"
	"log"
	"net"
	"net/rpc"
	"sort"
	"sync"
	"time"
)

func init() {
	gob.Register(time.Time{}) // Register time.Time for gob encoding/decoding
}

type BrokerService struct {
	config            *models.BrokerConfig
	storage           storage.Storage
	subscriberOffsets models.SubscriberOffsets
	mutex             sync.Mutex
	server            server.Server
}

func NewBrokerService(config *models.BrokerConfig) *BrokerService {
	broker := &BrokerService{}
	broker.config = config
	broker.storage = getStorage(config)
	broker.server = server.NewRpcServer()
	broker.subscriberOffsets = make(map[string]models.Offset)
	return broker
}

func (b *BrokerService) Start() {
	//register broker first
	if err := b.server.Register(b); err != nil {
		log.Fatalf("failed to register broker on rpc: %v\n", err)
	}

	//run on tcp
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.config.Host, b.config.Port))
	if err != nil {
		log.Println("failed to start broker:", err)
		return
	}
	defer listener.Close()
	log.Printf("Broker is running on %s:%s", b.config.Host, b.config.Port)

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
func (b *BrokerService) Publish(msg *models.Message, reply *string) error {
	if err := b.storage.StoreMessage(msg); err != nil {
		return err
	}

	*reply = fmt.Sprintf("message (%s) published successfully", msg.Id)
	return nil
}

// GetMessages retrieves messages for subscriber
func (b *BrokerService) GetMessages(request *models.GetMessageRequest, reply *models.MessageList) error {
	var err error
	switch request.PublishTime {
	case enums.Latest:
		*reply, err = b.getLatestMessages(request)
		return err
	case enums.Earliest:
		*reply, err = b.getEarliestMessages(request)
		return err
	//case enums.PreviousHour:
	//	*reply, err = b.getLatestMessages(request)
	//	return err
	//case enums.PreviousDay:
	//	*reply, err = b.getLatestMessages(request)
	//	return err
	default:
		*reply, err = b.getMessages(request, time.Duration(request.PublishTime))
		return err
	}
}

func (b *BrokerService) getEarliestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
	b.mutex.Lock()
	offset, exists := b.subscriberOffsets[request.SubscriberId]
	b.mutex.Unlock()
	//get messages from storage
	messages, err := b.storage.GetMessages(request.Topic, request.Partition)
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
		b.mutex.Lock()

		b.subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func (b *BrokerService) getLatestMessages(request *models.GetMessageRequest) (models.MessageList, error) {
	b.mutex.Lock()
	offset, exists := b.subscriberOffsets[request.SubscriberId]
	b.mutex.Unlock()
	//get messages from storage
	messages, err := b.storage.GetMessages(request.Topic, request.Partition)
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
		b.mutex.Lock()

		b.subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func (b *BrokerService) getMessages(request *models.GetMessageRequest, duration time.Duration) (models.MessageList, error) {
	b.mutex.Lock()
	offset, exists := b.subscriberOffsets[request.SubscriberId]
	b.mutex.Unlock()
	//get messages from storage
	messages, err := b.storage.GetMessages(request.Topic, request.Partition)
	if err != nil {
		return nil, err
	}

	if duration < 0 {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			return msg.Timestamp.Sub(time.Now().UTC().Add(duration)) >= 0
		})
	} else {
		messages = linq.Filter(messages, func(msg *models.Message, index int) bool {
			//d := msg.Timestamp.Sub(time.Now().UTC().Add(duration))
			//d := msg.Timestamp.Add(duration).Sub(msg.Timestamp)
			d := time.Now().Sub(msg.Timestamp.Add(duration))
			result := d <= 0
			return result
			return msg.Timestamp.Sub(time.Now().UTC().Add(duration)) <= 0
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
		b.mutex.Lock()

		b.subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: messages[len(messages)-1].Timestamp,
		}
		b.mutex.Unlock()
	} else {
		messages = []*models.Message{}
	}

	return messages, nil
}

func getStorage(config *models.BrokerConfig) storage.Storage {
	storageType := config.Storage
	switch storageType {
	case enums.StorageType_InMemory:
		return storage.NewInMemoryStorage()
	case enums.StorageType_Redis:
		return nil
	case enums.StorageType_Disk:
		return nil
	default:
		return storage.NewInMemoryStorage()
	}
}
