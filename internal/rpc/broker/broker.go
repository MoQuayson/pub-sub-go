package broker

import (
	"encoding/gob"
	"fmt"
	server "github.com/MoQuayson/pub-sub-go/internal/rpc"
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	models "github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/storage"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

func init() {
	gob.Register(time.Time{}) // Register time.Time for gob encoding/decoding
}

type BrokerService struct {
	broker.Broker
	config            *models.BrokerConfig
	storage           storage.Storage
	subscriberOffsets models.SubscriberOffsets
	mutex             sync.Mutex
	server            server.Server
}

func NewBrokerService(config *models.BrokerConfig) broker.Broker {
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
	b.mutex.Lock()
	if offset, exists := b.subscriberOffsets[request.SubscriberId]; exists {
		request.Timestamp = offset.Timestamp // Start from last read timestamp
	}
	b.mutex.Unlock()

	messages, err := b.storage.GetMessages(request.Topic, request.Partition)
	if err != nil {
		return err
	}

	var filteredMessages models.MessageList
	for _, msg := range messages {
		if msg.Timestamp.After(request.Timestamp) {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	if len(filteredMessages) > 0 {

		*reply = filteredMessages
		b.mutex.Lock()
		b.subscriberOffsets[request.SubscriberId] = models.Offset{
			Topic:     request.Topic,
			Partition: request.Partition,
			Timestamp: filteredMessages[len(filteredMessages)-1].Timestamp,
		}
		b.mutex.Unlock()
	} else {
		*reply = []*models.Message{}
	}

	return nil
}
