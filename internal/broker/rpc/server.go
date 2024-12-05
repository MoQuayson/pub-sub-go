package rpc

import (
	"fmt"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"log"
	"net/rpc"
)

//type Server interface {
//	Register() error
//	Publish(msg *models.Message, reply *string) error
//	GetMessages(req *models.GetMessageRequest, reply *models.MessageList) error
//}

type Server struct {
	broker *BrokerRpcServer
}

func NewRpcServer(br *BrokerRpcServer) *Server {
	return &Server{
		broker: br,
	}
}

func (r *Server) Register() error {
	if err := rpc.Register(r); err != nil {
		return err
	}

	return nil
}

func (r *Server) Publish(msg *models.Message, reply *string) error {
	if err := r.broker.PublishMessage(msg); err != nil {
		return err
	}

	*reply = fmt.Sprintf("message (%s) published successfully", msg.Id)
	return nil
}

func (r *Server) GetMessages(req *models.GetMessageRequest, reply *models.MessageList) error {
	//var err error
	messages, err := r.broker.GetMessages(req)
	if err != nil {
		log.Printf("failed to get messages: %v\n", err)
		return err
	}

	*reply = messages
	return nil
}
