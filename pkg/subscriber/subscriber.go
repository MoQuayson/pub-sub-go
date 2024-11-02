package subscriber

import (
	"github.com/MoQuayson/pub-sub-go/internal/subscriber/services"
	"github.com/MoQuayson/pub-sub-go/pkg/shared/models"
	"time"
)

type Subscriber interface {
	GetMessages(topic string, partition models.Partition, startTime time.Time) (models.MessageList, error)
}

func NewSubscriberService(cfg *models.RpcConnConfig) Subscriber {
	return services.NewSubscriberService(cfg)
}
