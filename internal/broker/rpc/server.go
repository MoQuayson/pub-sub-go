package rpc

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log/slog"
	"net/rpc"
	"strings"
)

type Server struct {
	broker *BrokerRpcServer
}

func NewRpcServer(br *BrokerRpcServer) *Server {
	return &Server{
		broker: br,
	}
}

func (r *Server) Register() error {
	err := rpc.Register(r)
	// net/rpc uses a global registry; ignore "already defined" when tests share a process.
	if err != nil && strings.Contains(err.Error(), "already defined") {
		return nil
	}
	return err
}

func (r *Server) Publish(msg *models.Message, reply *string) error {
	if err := r.broker.PublishMessage(msg); err != nil {
		return err
	}

	*reply = fmt.Sprintf("message (%s) published successfully", msg.Id)
	return nil
}

func (r *Server) GetMessages(req *models.GetMessageRequest, reply *models.MessageList) error {
	messages, err := r.broker.GetMessages(req)
	if err != nil {
		slog.Error("failed to get messages", "error", err)
		return err
	}

	*reply = messages
	return nil
}
