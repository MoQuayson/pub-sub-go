package rpc

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/MoQuayson/pub-sub-go/internal/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/storage"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log/slog"
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
	Logger            *slog.Logger
}

// Start runs the broker server, stopping cleanly when ctx is cancelled.
func (b *BrokerRpcServer) Start(ctx context.Context) error {
	if err := b.Server.Register(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", b.Config.Host, b.Config.Port))
	if err != nil {
		return fmt.Errorf("broker: listen on %s:%s: %w", b.Config.Host, b.Config.Port, err)
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		b.Logger.Info("broker shutting down")
		listener.Close()
	}()

	b.Logger.Info("broker listening", "host", b.Config.Host, "port", b.Config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("broker: accept: %w", err)
		}
		go rpc.ServeConn(conn)
	}
}

// PublishMessage publishes a message to a topic
func (b *BrokerRpcServer) PublishMessage(msg *models.Message) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

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
	return broker.ResolveMessages(&b.Mutex, b.SubscriberOffsets, messages, req)
}
