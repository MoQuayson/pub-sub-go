package broker

import "github.com/MoQuayson/pub-sub-go/pkg/shared/models"

type Broker interface {
	Start()
	Publish(msg *models.Message, reply *string) error
	GetMessages(request *models.GetMessageRequest, reply *models.MessageList) error
}
