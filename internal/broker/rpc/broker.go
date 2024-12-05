package rpc

import (
	"encoding/gob"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/internal/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

func init() {
	gob.Register(time.Time{}) // Register time.Time for gob encoding/decoding
}

type BrokerRpcServer struct {
	Config            *models.BrokerConfig
	Storage           storage.Storage
	SubscriberOffsets models.SubscriberOffsets
	Mutex             sync.Mutex
	Server            *Server
}

// Start  runs the broker server
func (b *BrokerRpcServer) Start() error {
	if err := b.Server.Register(); err != nil {
		return err
	}
	//run on tcp
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.Config.Host, b.Config.Port))
	if err != nil {
		log.Println("failed to start broker:", err)
		return err
	}
	defer listener.Close()
	log.Printf("Broker is running on %s:%s", b.Config.Host, b.Config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			return err
		}
		go rpc.ServeConn(conn)
	}

}

// PublishMessage publishes a message to a topic
func (b *BrokerRpcServer) PublishMessage(msg *models.Message) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	// Store the message
	if err := b.Storage.StoreMessage(msg); err != nil {
		return err
	}
	return nil
}

func (b *BrokerRpcServer) GetMessages(req *models.GetMessageRequest) ([]*models.Message, error) {
	messages, err := b.Storage.GetMessages(req.Topic, req.Partition)
	if err != nil {
		return nil, err
	}
	switch req.PublishTime {
	case models.LatestPublishTime:
		return broker.GetLatestMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
	case models.EarliestPublishTime:
		return broker.GetEarliestMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
	default:
		return broker.GetMessages(&b.Mutex, b.SubscriberOffsets, messages, req, time.Duration(req.PublishTime))
	}
}
